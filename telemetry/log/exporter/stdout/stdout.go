package stdout // import "go.microcore.dev/framework/telemetry/log/exporter/stdout"

import (
	"fmt"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
)

var logger = log.New(pkg)

func New(opts ...Option) (*stdoutlog.Exporter, error) {
	options := []stdoutlog.Option{}

	for _, opt := range opts {
		opt(&options)
	}

	exporter, err := stdoutlog.New(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	logger.Debug("exporter created")

	return exporter, nil
}
