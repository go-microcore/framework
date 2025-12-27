package periodic // import "go.microcore.dev/framework/telemetry/metric/reader/periodic"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/otel/sdk/metric"
)

var logger = log.New(pkg)

func New(exporter metric.Exporter, opts ...Option) *metric.PeriodicReader {
	options := []metric.PeriodicReaderOption{}
	for _, opt := range opts {
		opt(&options)
	}
	reader := metric.NewPeriodicReader(exporter, options...)
	logger.Debug("reader has been successfully created")
	return reader
}
