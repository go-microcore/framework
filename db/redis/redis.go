package redis // import "go.microcore.dev/framework/db/redis"

import (
	"context"
	"log/slog"
	"sync"
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/db/redis/client"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/telemetry"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type (
	Manager interface {
		Client() *redis.Client
		SetClient(client *redis.Client) Manager
		SetTelemetryManager(telemetry telemetry.Manager) error
		GetShutdownTimeout() time.Duration
		GetShutdownHandler() bool
		Shutdown(ctx context.Context, reason string) error
	}

	r struct {
		client          *redis.Client
		shutdownTimeout time.Duration
		shutdownHandler bool
		mu              sync.RWMutex
	}
)

// Nil reply returned by Redis when key does not exist.
const Nil = redis.Nil

var logger = log.New(pkg)

func New(opts ...Option) Manager {
	r := &r{
		shutdownTimeout: DefaultShutdownTimeout,
		shutdownHandler: DefaultShutdownHandler,
	}

	for _, opt := range opts {
		opt(r)
	}

	if r.client == nil {
		r.client = client.New()
	}

	if r.shutdownHandler {
		shutdown.AddHandler(r.Shutdown)
		logger.Debug("shutdown handler has been successfully registered")
	}

	logger.Info(
		"manager has been successfully created",
		slog.Group("shutdown",
			slog.Duration("timeout", r.shutdownTimeout),
			slog.Bool("handler", r.shutdownHandler),
		),
	)

	return r
}

func (r *r) Client() *redis.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.client
}

func (r *r) SetClient(client *redis.Client) Manager {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.client = client
	return r
}

func (r *r) SetTelemetryManager(telemetry telemetry.Manager) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := redisotel.InstrumentTracing(
		r.client,
		redisotel.WithTracerProvider(
			telemetry.GetTraceProvider(),
		),
	); err != nil {
		logger.Error(
			"failed to use otel tracing",
			slog.Any("error", err),
		)
		return err
	}

	logger.Info("otel tracing has been successfully initialized")
	return nil
}

func (r *r) GetShutdownTimeout() time.Duration {
	return r.shutdownTimeout
}

func (r *r) GetShutdownHandler() bool {
	return r.shutdownHandler
}

func (r *r) Shutdown(ctx context.Context, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, r.shutdownTimeout)
	defer cancel()

	logger.Debug(
		"shutdown",
		slog.String("reason", reason),
	)

	ch := make(chan error, 1)
	go func() {
		ch <- r.client.Close()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}
