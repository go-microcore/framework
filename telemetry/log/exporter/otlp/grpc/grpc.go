package grpc // import "go.microcore.dev/framework/telemetry/log/exporter/otlp/grpc"

import (
	"context"
	"fmt"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
)

var logger = log.New(pkg)

func New(ctx context.Context, opts ...Option) (*otlploggrpc.Exporter, error) {
	options := []otlploggrpc.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := otlploggrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	logger.Debug("exporter created")

	return exporter, nil
}
