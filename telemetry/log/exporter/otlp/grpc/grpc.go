package grpc // import "go.microcore.dev/framework/telemetry/log/exporter/otlp/grpc"

import (
	"context"
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
)

var logger = log.New(pkg)

func New(ctx context.Context, opts ...Option) *otlploggrpc.Exporter {
	options := []otlploggrpc.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := otlploggrpc.New(ctx, options...)
	if err != nil {
		logger.Error(
			"failed to create exporter",
			slog.Any("error", err),
		)
		panic(err)
	}

	logger.Info("exporter has been successfully created")

	return exporter
}
