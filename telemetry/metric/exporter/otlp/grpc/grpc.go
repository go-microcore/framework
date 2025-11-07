package grpc // import "go.microcore.dev/framework/telemetry/metric/exporter/otlp/grpc"

import (
	"context"
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
)

var logger = log.New(pkg)

func New(ctx context.Context, opts ...Option) *otlpmetricgrpc.Exporter {
	options := []otlpmetricgrpc.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := otlpmetricgrpc.New(ctx, options...)
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
