package stdout // import "go.microcore.dev/framework/telemetry/log/exporter/stdout"

import (
	"io"

	_ "go.microcore.dev/framework"

	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
)

type Option func(*[]stdoutlog.Option)

// WithWriter sets the export stream destination.
func WithWriter(writer io.Writer) Option {
	return func(o *[]stdoutlog.Option) {
		*o = append(*o, stdoutlog.WithWriter(writer))
	}
}

// WithPrettyPrint prettifies the emitted output.
func WithPrettyPrint() Option {
	return func(o *[]stdoutlog.Option) {
		*o = append(*o, stdoutlog.WithPrettyPrint())
	}
}

// WithoutTimestamps sets the export stream to not include timestamps.
func WithoutTimestamps() Option {
	return func(o *[]stdoutlog.Option) {
		*o = append(*o, stdoutlog.WithoutTimestamps())
	}
}
