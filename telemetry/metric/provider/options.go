package provider // import "go.microcore.dev/framework/telemetry/metric/provider"

import (
	_ "go.microcore.dev/framework"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Option func(*[]metric.Option)

// WithResource associates a Resource with a MeterProvider. This Resource
// represents the entity producing telemetry and is associated with all Meters
// the MeterProvider will create.
//
// By default, if this Option is not used, the default Resource from the
// go.opentelemetry.io/otel/sdk/resource package will be used.
func WithResource(res *resource.Resource) Option {
	return func(o *[]metric.Option) {
		*o = append(*o, metric.WithResource(res))
	}
}

// WithReader associates Reader r with a MeterProvider.
//
// By default, if this option is not used, the MeterProvider will perform no
// operations; no data will be exported without a Reader.
func WithReader(r metric.Reader) Option {
	return func(o *[]metric.Option) {
		*o = append(*o, metric.WithReader(r))
	}
}

// WithView associates views with a MeterProvider.
//
// Views are appended to existing ones in a MeterProvider if this option is
// used multiple times.
//
// By default, if this option is not used, the MeterProvider will use the
// default view.
func WithView(views ...metric.View) Option {
	return func(o *[]metric.Option) {
		*o = append(*o, metric.WithView(views...))
	}
}

// WithExemplarFilter configures the exemplar filter.
//
// The exemplar filter determines which measurements are offered to the
// exemplar reservoir, but the exemplar reservoir makes the final decision of
// whether to store an exemplar.
//
// By default, the [exemplar.SampledFilter]
// is used. Exemplars can be entirely disabled by providing the
// [exemplar.AlwaysOffFilter].
func WithExemplarFilter(filter exemplar.Filter) Option {
	return func(o *[]metric.Option) {
		*o = append(*o, metric.WithExemplarFilter(filter))
	}
}

// WithCardinalityLimit sets the cardinality limit for the MeterProvider.
//
// The cardinality limit is the hard limit on the number of metric datapoints
// that can be collected for a single instrument in a single collect cycle.
//
// Setting this to a zero or negative value means no limit is applied.
func WithCardinalityLimit(limit int) Option {
	return func(o *[]metric.Option) {
		*o = append(*o, metric.WithCardinalityLimit(limit))
	}
}
