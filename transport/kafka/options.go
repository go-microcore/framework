package kafka // import "go.microcore.dev/framework/transport/kafka"

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/errors"
	"go.microcore.dev/framework/telemetry"
)

type Option func(*k)

func WithWriterBrokers(brokers []string) Option {
	return func(k *k) {
		k.brokers.writer = brokers
	}
}

func WithReaderBrokers(brokers []string) Option {
	return func(k *k) {
		k.brokers.reader = brokers
	}
}

func WithWriters(writers map[string]*kafka.Writer) Option {
	return func(k *k) {
		k.writers = writers
	}
}

func WithReaders(readers map[string]*kafka.Reader) Option {
	return func(k *k) {
		k.readers = readers
	}
}

func WithTelemetryManager(telemetry telemetry.Manager) Option {
	return func(k *k) {
		k.telemetry = telemetry
	}
}

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(k *k) {
		k.shutdownTimeout = timeout
	}
}

func WithoutShutdownHandler() Option {
	return func(k *k) {
		k.shutdownHandler = false
	}
}

type PubOption func(*pub)

func WithPubContext(context context.Context) PubOption {
	return func(p *pub) {
		p.context = context
	}
}

func WithPubPartition(partition int) PubOption {
	return func(p *pub) {
		p.message.Partition = partition
	}
}

func WithPubOffset(offset int64) PubOption {
	return func(p *pub) {
		p.message.Offset = offset
	}
}

func WithPubHighWaterMark(highWaterMark int64) PubOption {
	return func(p *pub) {
		p.message.HighWaterMark = highWaterMark
	}
}

func WithPubKey(key []byte) PubOption {
	return func(p *pub) {
		p.message.Key = key
	}
}

func WithPubHeader(header kafka.Header) PubOption {
	return func(p *pub) {
		p.message.Headers = append(p.message.Headers, header)
	}
}

// This field is used to hold arbitrary data you wish to include, so it
// will be available when handle it on the Writer's `Completion` method,
// this support the application can do any post operation on each message.
func WithPubWriterData(writerData any) PubOption {
	return func(p *pub) {
		p.message.WriterData = writerData
	}
}

// If not set at the creation, Time will be automatically set when
// writing the message.
func WithPubTime(time time.Time) PubOption {
	return func(p *pub) {
		p.message.Time = time
	}
}

type SubOption func(*sub)

func WithSubContext(context context.Context) SubOption {
	return func(s *sub) {
		s.context = context
	}
}

func WithSubHandler(handler func(ctx context.Context, message kafka.Message) error) SubOption {
	return func(s *sub) {
		s.handler = handler
	}
}

func WithSubPayloadParserHandler[T any](handler func(ctx context.Context, message kafka.Message, payload *T) error) SubOption {
	return func(s *sub) {
		s.handler = func(ctx context.Context, message kafka.Message) error {
			var payload T
			if err := json.Unmarshal(message.Value, &payload); err == nil {
				if v, ok := any(&payload).(interface{ Validate() error }); ok {
					if err := v.Validate(); err != nil {
						return err
					}
				}
			} else {
				return errors.ErrUnsupportedMediaType
			}
			return handler(ctx, message, &payload)
		}
	}
}
