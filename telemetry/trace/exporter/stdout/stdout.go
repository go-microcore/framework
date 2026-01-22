package stdout // import "go.microcore.dev/framework/telemetry/trace/exporter/stdout"

import (
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

var logger = log.New(pkg)

func New(opts ...Option) trace.SpanExporter {
	options := []stdouttrace.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := stdouttrace.New(options...)
	if err != nil {
		logger.Error(
			"failed to create exporter",
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitUnavailable)
	}

	logger.Debug("exporter has been successfully created")

	return exporter
}
