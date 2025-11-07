package telemetry // import "go.microcore.dev/framework/telemetry"

import (
	"time"

	_ "go.microcore.dev/framework"
	logProvider "go.microcore.dev/framework/telemetry/log/provider"
	metricProvider "go.microcore.dev/framework/telemetry/metric/provider"
	traceProvider "go.microcore.dev/framework/telemetry/trace/provider"

	"go.opentelemetry.io/otel/propagation"
	logSdk "go.opentelemetry.io/otel/sdk/log"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

type Option func(*t)

func WithTraceProvider(provider *traceSdk.TracerProvider) Option {
	return func(t *t) {
		t.traceProvider = provider
	}
}

func WithTraceProviderOptions(opts ...traceProvider.Option) Option {
	return func(t *t) {
		t.traceProvider = traceProvider.New(opts...)
	}
}

func WithMetricProvider(provider *metricSdk.MeterProvider) Option {
	return func(t *t) {
		t.metricProvider = provider
	}
}

func WithMetricProviderOptions(opts ...metricProvider.Option) Option {
	return func(t *t) {
		t.metricProvider = metricProvider.New(opts...)
	}
}

func WithLogProvider(provider *logSdk.LoggerProvider) Option {
	return func(t *t) {
		t.logProvider = provider
	}
}

func WithLogProviderOptions(opts ...logProvider.Option) Option {
	return func(t *t) {
		t.logProvider = logProvider.New(opts...)
	}
}

func WithPropagator(propagator propagation.TextMapPropagator) Option {
	return func(t *t) {
		t.propagator = propagator
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(t *t) {
		t.shutdownTimeout = timeout
	}
}

func WithoutShutdownHandler() Option {
	return func(t *t) {
		t.shutdownHandler = false
	}
}

func WithoutSetLogProvider() Option {
	return func(t *t) {
		t.setLogProvider = false
	}
}
