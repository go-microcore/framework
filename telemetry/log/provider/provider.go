package provider // import "go.microcore.dev/framework/telemetry/log/provider"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	logSdk "go.opentelemetry.io/otel/sdk/log"
)

var logger = log.New(pkg)

func New(opts ...Option) *logSdk.LoggerProvider {
	options := []logSdk.LoggerProviderOption{}
	for _, opt := range opts {
		opt(&options)
	}
	provider := logSdk.NewLoggerProvider(options...)
	logger.Debug("provider created")
	return provider
}
