package grpc // import "go.microcore.dev/framework/telemetry/trace/exporter/otlp/grpc"

import (
	"context"
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

var logger = log.New(pkg)

func New(ctx context.Context, opts ...Option) *otlptrace.Exporter {
	options := []otlptracegrpc.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := otlptracegrpc.New(ctx, options...)
	if err != nil {
		logger.Error(
			"failed to create exporter",
			slog.Any("error", err),
		)
		panic(err)
	}

	logger.Debug("exporter has been successfully created")

	return exporter
}
