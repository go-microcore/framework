package periodic // import "go.microcore.dev/framework/telemetry/metric/reader/periodic"

import (
	"time"

	_ "go.microcore.dev/framework"

	"go.opentelemetry.io/otel/sdk/metric"
)

type Option func(*[]metric.PeriodicReaderOption)

// WithTimeout configures the time a PeriodicReader waits for an export to
// complete before canceling it. This includes an export which occurs as part
// of Shutdown or ForceFlush if the user passed context does not have a
// deadline. If the user passed context does have a deadline, it will be used
// instead.
//
// This option overrides any value set for the
// OTEL_METRIC_EXPORT_TIMEOUT environment variable.
//
// If this option is not used or d is less than or equal to zero, 30 seconds
// is used as the default.
func WithTimeout(d time.Duration) Option {
	return func(o *[]metric.PeriodicReaderOption) {
		*o = append(*o, metric.WithTimeout(d))
	}
}

// WithInterval configures the intervening time between exports for a
// PeriodicReader.
//
// This option overrides any value set for the
// OTEL_METRIC_EXPORT_INTERVAL environment variable.
//
// If this option is not used or d is less than or equal to zero, 60 seconds
// is used as the default.
func WithInterval(d time.Duration) Option {
	return func(o *[]metric.PeriodicReaderOption) {
		*o = append(*o, metric.WithInterval(d))
	}
}
