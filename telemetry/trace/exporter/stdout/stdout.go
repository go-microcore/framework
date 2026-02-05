package stdout // import "go.microcore.dev/framework/telemetry/trace/exporter/stdout"

import (
	"fmt"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

var logger = log.New(pkg)

func New(opts ...Option) (trace.SpanExporter, error) {
	options := []stdouttrace.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := stdouttrace.New(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	logger.Debug("exporter created")

	return exporter, nil
}
