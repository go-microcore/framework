package log // import "go.microcore.dev/framework/telemetry/log"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"go.opentelemetry.io/contrib/processors/minsev"
	otelLog "go.opentelemetry.io/otel/log"
	logSdk "go.opentelemetry.io/otel/sdk/log"
)

type severity struct{}

func (s severity) Severity() otelLog.Severity {
	var sev minsev.Severity
	_ = sev.UnmarshalText([]byte(log.Level().String()))
	return sev.Severity()
}

func NewProcessor(processor logSdk.Processor) *minsev.LogProcessor {
	return minsev.NewLogProcessor(
		processor,
		severity{},
	)
}
