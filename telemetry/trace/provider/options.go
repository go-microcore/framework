package provider // import "go.microcore.dev/framework/telemetry/trace/provider"

import (
	_ "go.microcore.dev/framework"
	
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Option func(*[]trace.TracerProviderOption)

// WithSyncer registers the exporter with the TracerProvider using a
// SimpleSpanProcessor.
//
// This is not recommended for production use. The synchronous nature of the
// SimpleSpanProcessor that will wrap the exporter make it good for testing,
// debugging, or showing examples of other feature, but it will be slow and
// have a high computation resource usage overhead. The WithBatcher option is
// recommended for production use instead.
func WithSyncer(e trace.SpanExporter) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithSyncer(e))
	}
}

// WithBatcher registers the exporter with the TracerProvider using a
// BatchSpanProcessor configured with the passed opts.
func WithBatcher(e trace.SpanExporter, opts ...trace.BatchSpanProcessorOption) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithBatcher(e, opts...))
	}
}

// WithSpanProcessor registers the SpanProcessor with a TracerProvider.
func WithSpanProcessor(sp trace.SpanProcessor) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithSpanProcessor(sp))
	}
}

// WithResource returns a TracerProviderOption that will configure the
// Resource r as a TracerProvider's Resource. The configured Resource is
// referenced by all the Tracers the TracerProvider creates. It represents the
// entity producing telemetry.
//
// If this option is not used, the TracerProvider will use the
// resource.Default() Resource by default.
func WithResource(r *resource.Resource) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithResource(r))
	}
}

// WithIDGenerator returns a TracerProviderOption that will configure the
// IDGenerator g as a TracerProvider's IDGenerator. The configured IDGenerator
// is used by the Tracers the TracerProvider creates to generate new Span and
// Trace IDs.
//
// If this option is not used, the TracerProvider will use a random number
// IDGenerator by default.
func WithIDGenerator(g trace.IDGenerator) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithIDGenerator(g))
	}
}

// WithSampler returns a TracerProviderOption that will configure the Sampler
// s as a TracerProvider's Sampler. The configured Sampler is used by the
// Tracers the TracerProvider creates to make their sampling decisions for the
// Spans they create.
//
// This option overrides the Sampler configured through the OTEL_TRACES_SAMPLER
// and OTEL_TRACES_SAMPLER_ARG environment variables. If this option is not used
// and the sampler is not configured through environment variables or the environment
// contains invalid/unsupported configuration, the TracerProvider will use a
// ParentBased(AlwaysSample) Sampler by default.
func WithSampler(s trace.Sampler) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithSampler(s))
	}
}

// WithSpanLimits returns a TracerProviderOption that configures a
// TracerProvider to use the SpanLimits sl. These SpanLimits bound any Span
// created by a Tracer from the TracerProvider.
//
// If any field of sl is zero or negative it will be replaced with the default
// value for that field.
//
// If this or WithRawSpanLimits are not provided, the TracerProvider will use
// the limits defined by environment variables, or the defaults if unset.
// Refer to the NewSpanLimits documentation for information about this
// relationship.
//
// Deprecated: Use WithRawSpanLimits instead which allows setting unlimited
// and zero limits. This option will be kept until the next major version
// incremented release.
func WithSpanLimits(sl trace.SpanLimits) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithSpanLimits(sl))
	}
}

// WithRawSpanLimits returns a TracerProviderOption that configures a
// TracerProvider to use these limits. These limits bound any Span created by
// a Tracer from the TracerProvider.
//
// The limits will be used as-is. Zero or negative values will not be changed
// to the default value like WithSpanLimits does. Setting a limit to zero will
// effectively disable the related resource it limits and setting to a
// negative value will mean that resource is unlimited. Consequentially, this
// means that the zero-value SpanLimits will disable all span resources.
// Because of this, limits should be constructed using NewSpanLimits and
// updated accordingly.
//
// If this or WithSpanLimits are not provided, the TracerProvider will use the
// limits defined by environment variables, or the defaults if unset. Refer to
// the NewSpanLimits documentation for information about this relationship.
func WithRawSpanLimits(limits trace.SpanLimits) Option {
	return func(o *[]trace.TracerProviderOption) {
		*o = append(*o, trace.WithRawSpanLimits(limits))
	}
}
