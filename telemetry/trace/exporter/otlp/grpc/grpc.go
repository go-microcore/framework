package grpc // import "go.microcore.dev/framework/telemetry/trace/exporter/otlp/grpc"

import (
	"context"
	"fmt"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

var logger = log.New(pkg)

func New(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	options := []otlptracegrpc.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	logger.Debug("exporter created")

	return exporter, nil
}
