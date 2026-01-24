package provider // import "go.microcore.dev/framework/telemetry/trace/provider"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/sdk/trace"
)

var logger = log.New(pkg)

func New(opts ...Option) *trace.TracerProvider {
	options := []trace.TracerProviderOption{}

	for _, opt := range opts {
		opt(&options)
	}
	provider := trace.NewTracerProvider(options...)
	logger.Debug("provider created")
	return provider
}
