package stdout // import "go.microcore.dev/framework/telemetry/log/exporter/stdout"

import (
	"log/slog"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
)

var logger = log.New(pkg)

func New(opts ...Option) *stdoutlog.Exporter {
	options := []stdoutlog.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := stdoutlog.New(options...)
	if err != nil {
		logger.Error(
			"failed to create exporter",
			slog.Any("error", err),
		)
		shutdown.Exit(shutdown.ExitUnavailable)
	}

	logger.Debug("exporter created")

	return exporter
}
