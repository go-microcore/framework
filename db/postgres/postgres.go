package postgres // import "go.microcore.dev/framework/db/postgres"

import (
	"context"
	"log/slog"
	"sync"
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/db/postgres/client"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/telemetry"

	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/go-gormigrate/gormigrate/v2"
)

type (
	Manager interface {
		Client() *gorm.DB
		SetClient(client *gorm.DB) Manager
		SetTelemetryManager(telemetry telemetry.Manager) error
		Migrate(migrations []*gormigrate.Migration, options *gormigrate.Options) error
		GetShutdownTimeout() time.Duration
		GetShutdownHandler() bool
		Shutdown(ctx context.Context, code int) error
	}

	p struct {
		client          *gorm.DB
		shutdownTimeout time.Duration
		shutdownHandler bool
		mu              sync.RWMutex
	}
)

var logger = log.New(pkg)

func New(opts ...Option) Manager {
	p := &p{
		shutdownTimeout: DefaultShutdownTimeout,
		shutdownHandler: DefaultShutdownHandler,
	}

	for _, opt := range opts {
		opt(p)
	}

	if p.client == nil {
		p.client = client.New()
	}

	if p.shutdownHandler {
		shutdown.AddHandler(p.Shutdown)
		logger.Debug("shutdown handler has been successfully registered")
	}

	logger.Info(
		"manager has been successfully created",
		slog.Group("shutdown",
			slog.Duration("timeout", p.shutdownTimeout),
			slog.Bool("handler", p.shutdownHandler),
		),
	)

	return p
}

func (p *p) Client() *gorm.DB {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.client
}

func (p *p) SetClient(client *gorm.DB) Manager {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.client = client
	return p
}

func (p *p) SetTelemetryManager(telemetry telemetry.Manager) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.client.Use(
		tracing.NewPlugin(
			tracing.WithTracerProvider(
				telemetry.GetTraceProvider(),
			),
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

func (p *p) Migrate(migrations []*gormigrate.Migration, options *gormigrate.Options) error {
	m := gormigrate.New(p.client, options, migrations)
	if err := m.Migrate(); err != nil {
		logger.Error(
			"could not migrate",
			slog.Any("error", err),
		)
		return err
	}
	logger.Info("migrations have been successfully completed")
	return nil
}

func (p *p) GetShutdownTimeout() time.Duration {
	return p.shutdownTimeout
}

func (p *p) GetShutdownHandler() bool {
	return p.shutdownHandler
}

func (p *p) Shutdown(ctx context.Context, code int) error {
	ctx, cancel := context.WithTimeout(ctx, p.shutdownTimeout)
	defer cancel()

	logger.Debug(
		"shutdown",
		slog.Int("code", code),
	)

	p.mu.RLock()
	db, err := p.client.DB()
	p.mu.RUnlock()
	if err != nil {
		return err
	}

	ch := make(chan error, 1)
	go func() {
		ch <- db.Close()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}
