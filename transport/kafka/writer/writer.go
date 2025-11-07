package writer // import "go.microcore.dev/framework/transport/kafka/writer"

import (
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"

	"github.com/segmentio/kafka-go"
)

var logger = log.New(pkg)

func New(opts ...Option) *kafka.Writer {
	writer := &kafka.Writer{}
	for _, opt := range opts {
		opt(writer)
	}
	logger.Info("writer has been successfully created")
	return writer
}
