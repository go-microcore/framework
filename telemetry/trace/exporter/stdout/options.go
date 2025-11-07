package stdout // import "go.microcore.dev/framework/telemetry/trace/exporter/stdout"

import (
	"io"

	_ "go.microcore.dev/framework"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

type Option func(*[]stdouttrace.Option)

// WithWriter sets the export stream destination.
func WithWriter(writer io.Writer) Option {
	return func(o *[]stdouttrace.Option) {
		*o = append(*o, stdouttrace.WithWriter(writer))
	}
}

// WithPrettyPrint prettifies the emitted output.
func WithPrettyPrint() Option {
	return func(o *[]stdouttrace.Option) {
		*o = append(*o, stdouttrace.WithPrettyPrint())
	}
}

// WithoutTimestamps sets the export stream to not include timestamps.
func WithoutTimestamps() Option {
	return func(o *[]stdouttrace.Option) {
		*o = append(*o, stdouttrace.WithoutTimestamps())
	}
}
