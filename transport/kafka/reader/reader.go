package reader // import "go.microcore.dev/framework/transport/kafka/reader"

import (
	"github.com/segmentio/kafka-go"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
)

var logger = log.New(pkg)

func New(opts ...Option) *kafka.Reader {
	config := &kafka.ReaderConfig{}
	for _, opt := range opts {
		opt(config)
	}
	reader := kafka.NewReader(*config)
	logger.Debug("reader created")
	return reader
}
