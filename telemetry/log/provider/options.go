package provider // import "go.microcore.dev/framework/telemetry/log/provider"

import (
	_ "go.microcore.dev/framework"
	
	logSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

type Option func(*[]logSdk.LoggerProviderOption)

// WithResource associates a Resource with a LoggerProvider. This Resource
// represents the entity producing telemetry and is associated with all Loggers
// the LoggerProvider will create.
//
// By default, if this Option is not used, the default Resource from the
// go.opentelemetry.io/otel/sdk/resource package will be used.
func WithResource(res *resource.Resource) Option {
	return func(o *[]logSdk.LoggerProviderOption) {
		*o = append(*o, logSdk.WithResource(res))
	}
}

// WithProcessor associates Processor with a LoggerProvider.
//
// By default, if this option is not used, the LoggerProvider will perform no
// operations; no data will be exported without a processor.
//
// The SDK invokes the processors sequentially in the same order as they were
// registered.
//
// For production, use [NewBatchProcessor] to batch log records before they are exported.
// For testing and debugging, use [NewSimpleProcessor] to synchronously export log records.
//
// See [FilterProcessor] for information about how a Processor can support filtering.
func WithProcessor(processor logSdk.Processor) Option {
	return func(o *[]logSdk.LoggerProviderOption) {
		*o = append(*o, logSdk.WithProcessor(processor))
	}
}

// WithAttributeCountLimit sets the maximum allowed log record attribute count.
// Any attribute added to a log record once this limit is reached will be dropped.
//
// Setting this to zero means no attributes will be recorded.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_COUNT_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, 128 will be used.
func WithAttributeCountLimit(limit int) Option {
	return func(o *[]logSdk.LoggerProviderOption) {
		*o = append(*o, logSdk.WithAttributeCountLimit(limit))
	}
}

// WithAttributeValueLengthLimit sets the maximum allowed attribute value length.
//
// This limit only applies to string and string slice attribute values.
// Any string longer than this value will be truncated to this length.
//
// Setting this to a negative value means no limit is applied.
//
// If the OTEL_LOGRECORD_ATTRIBUTE_VALUE_LENGTH_LIMIT environment variable is set,
// and this option is not passed, that variable value will be used.
//
// By default, if an environment variable is not set, and this option is not
// passed, no limit (-1) will be used.
func WithAttributeValueLengthLimit(limit int) Option {
	return func(o *[]logSdk.LoggerProviderOption) {
		*o = append(*o, logSdk.WithAttributeValueLengthLimit(limit))
	}
}

// WithAllowKeyDuplication sets whether deduplication is skipped for log attributes or other key-value collections.
//
// By default, the key-value collections within a log record are deduplicated to comply with the OpenTelemetry Specification.
// Deduplication means that if multiple key–value pairs with the same key are present, only a single pair
// is retained and others are discarded.
//
// Disabling deduplication with this option can improve performance e.g. of adding attributes to the log record.
//
// Note that if you disable deduplication, you are responsible for ensuring that duplicate
// key-value pairs within in a single collection are not emitted,
// or that the telemetry receiver can handle such duplicates.
func WithAllowKeyDuplication() Option {
	return func(o *[]logSdk.LoggerProviderOption) {
		*o = append(*o, logSdk.WithAllowKeyDuplication())
	}
}
