package telemetry // import "go.microcore.dev/framework/telemetry"

import (
	"time"

	_ "go.microcore.dev/framework"
	"go.opentelemetry.io/otel/propagation"
)

const (
	pkg = "go.microcore.dev/framework/telemetry"

	DefaultShutdownTimeout = 10 * time.Second
	DefaultShutdownHandler = true
	DefaultSetLogProvider  = true

	DefaultMetricPeriodicReaderInterval = 30 * time.Second
	DefaultMetricPeriodicReaderTimeout  = 10 * time.Second

	DefaultLogExportInterval = 30 * time.Second
	DefaultLogExportTimeout  = 10 * time.Second
)

var (
	defaultPropagator = &propagation.TraceContext{}
)
