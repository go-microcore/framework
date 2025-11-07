package provider // import "go.microcore.dev/framework/telemetry/metric/provider"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/sdk/metric"
)

var logger = log.New(pkg)

func New(opts ...Option) *metric.MeterProvider {
	options := []metric.Option{}
	for _, opt := range opts {
		opt(&options)
	}
	provider := metric.NewMeterProvider(options...)
	logger.Info("provider has been successfully created")
	return provider
}
