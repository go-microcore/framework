package redis // import "go.microcore.dev/framework/db/redis"

import (
	"context"
	"log/slog"
	"os"
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

type Interface interface {
	Client() *redis.Client
	SetClient(client *redis.Client) Interface
	UseTelemetry(telemetry telemetry.Interface) error
	GetShutdownTimeout() time.Duration
	GetShutdownHandler() bool
	Shutdown(ctx context.Context, reason string) error
	ShutdownHandler(sig os.Signal) error
}

type r struct {
	client          *redis.Client
	shutdownTimeout time.Duration
	shutdownHandler bool
	mu              sync.RWMutex
}

var logger = log.New(pkg)

func New(opts ...Option) Interface {
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
		shutdown.AddHandler(r.ShutdownHandler)
		logger.Info("shutdown handler has been successfully registered")
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

func (r *r) SetClient(client *redis.Client) Interface {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.client = client
	return r
}

func (r *r) UseTelemetry(telemetry telemetry.Interface) error {
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
	logger.Info(
		"shutting down",
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

func (r *r) ShutdownHandler(sig os.Signal) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.shutdownTimeout)
	defer cancel()

	reason := "unknown"
	if sig != nil {
		reason = sig.String()
	}

	return r.Shutdown(ctx, reason)
}
