package writer // import "go.microcore.dev/framework/transport/kafka/writer"

import (
	"net"
	"time"

	_ "go.microcore.dev/framework"

	"github.com/segmentio/kafka-go"
)

type Option func(*kafka.Writer)

// Address of the kafka cluster that this writer is configured to send
// messages to.
//
// This field is required, attempting to write messages to a writer with a
// nil address will error.
func WithAddr(addr net.Addr) Option {
	return func(k *kafka.Writer) {
		k.Addr = addr
	}
}

// Topic is the name of the topic that the writer will produce messages to.
//
// Setting this field or not is a mutually exclusive option. If you set Topic
// here, you must not set Topic for any produced Message. Otherwise, if you	do
// not set Topic, every Message must have Topic specified.
func WithTopic(topic string) Option {
	return func(k *kafka.Writer) {
		k.Topic = topic
	}
}

// The balancer used to distribute messages across partitions.
//
// The default is to use a round-robin distribution.
func WithBalancer(balancer kafka.Balancer) Option {
	return func(k *kafka.Writer) {
		k.Balancer = balancer
	}
}

// Limit on how many attempts will be made to deliver a message.
//
// The default is to try at most 10 times.
func WithMaxAttempts(maxAttempts int) Option {
	return func(k *kafka.Writer) {
		k.MaxAttempts = maxAttempts
	}
}

// WriteBackoffMin optionally sets the smallest amount of time the writer waits before
// it attempts to write a batch of messages
//
// Default: 100ms
func WithWriteBackoffMin(writeBackoffMin time.Duration) Option {
	return func(k *kafka.Writer) {
		k.WriteBackoffMin = writeBackoffMin
	}
}

// WriteBackoffMax optionally sets the maximum amount of time the writer waits before
// it attempts to write a batch of messages
//
// Default: 1s
func WithWriteBackoffMax(writeBackoffMax time.Duration) Option {
	return func(k *kafka.Writer) {
		k.WriteBackoffMax = writeBackoffMax
	}
}

// Limit on how many messages will be buffered before being sent to a
// partition.
//
// The default is to use a target batch size of 100 messages.
func WithBatchSize(batchSize int) Option {
	return func(k *kafka.Writer) {
		k.BatchSize = batchSize
	}
}

// Limit the maximum size of a request in bytes before being sent to
// a partition.
//
// The default is to use a kafka default value of 1048576.
func WithBatchBytes(batchBytes int64) Option {
	return func(k *kafka.Writer) {
		k.BatchBytes = batchBytes
	}
}

// Time limit on how often incomplete message batches will be flushed to
// kafka.
//
// The default is to flush at least every second.
func WithBatchTimeout(batchTimeout time.Duration) Option {
	return func(k *kafka.Writer) {
		k.BatchTimeout = batchTimeout
	}
}

// Timeout for read operations performed by the Writer.
//
// Defaults to 10 seconds.
func WithReadTimeout(readTimeout time.Duration) Option {
	return func(k *kafka.Writer) {
		k.ReadTimeout = readTimeout
	}
}

// Timeout for write operation performed by the Writer.
//
// Defaults to 10 seconds.
func WithWriteTimeout(writeTimeout time.Duration) Option {
	return func(k *kafka.Writer) {
		k.WriteTimeout = writeTimeout
	}
}

// Number of acknowledges from partition replicas required before receiving
// a response to a produce request, the following values are supported:
//
//	RequireNone (0)  fire-and-forget, do not wait for acknowledgements from the
//	RequireOne  (1)  wait for the leader to acknowledge the writes
//	RequireAll  (-1) wait for the full ISR to acknowledge the writes
//
// Defaults to RequireNone.
func WithRequiredAcks(requiredAcks kafka.RequiredAcks) Option {
	return func(k *kafka.Writer) {
		k.RequiredAcks = requiredAcks
	}
}

// Setting this flag to true causes the WriteMessages method to never block.
// It also means that errors are ignored since the caller will not receive
// the returned value. Use this only if you don't care about guarantees of
// whether the messages were written to kafka.
//
// Defaults to false.
func WithAsync(async bool) Option {
	return func(k *kafka.Writer) {
		k.Async = async
	}
}

// An optional function called when the writer succeeds or fails the
// delivery of messages to a kafka partition. When writing the messages
// fails, the `err` parameter will be non-nil.
//
// The messages that the Completion function is called with have their
// topic, partition, offset, and time set based on the Produce responses
// received from kafka. All messages passed to a call to the function have
// been written to the same partition. The keys and values of messages are
// referencing the original byte slices carried by messages in the calls to
// WriteMessages.
//
// The function is called from goroutines started by the writer. Calls to
// Close will block on the Completion function calls. When the Writer is
// not writing asynchronously, the WriteMessages call will also block on
// Completion function, which is a useful guarantee if the byte slices
// for the message keys and values are intended to be reused after the
// WriteMessages call returned.
//
// If a completion function panics, the program terminates because the
// panic is not recovered by the writer and bubbles up to the top of the
// goroutine's call stack.
func WithCompletion(completion func(messages []kafka.Message, err error)) Option {
	return func(k *kafka.Writer) {
		k.Completion = completion
	}
}

// Compression set the compression codec to be used to compress messages.
func WithCompression(compression kafka.Compression) Option {
	return func(k *kafka.Writer) {
		k.Compression = compression
	}
}

// If not nil, specifies a logger used to report internal changes within the
// writer.
func WithLogger(logger kafka.Logger) Option {
	return func(k *kafka.Writer) {
		k.Logger = logger
	}
}

// ErrorLogger is the logger used to report errors. If nil, the writer falls
// back to using Logger instead.
func WithErrorLogger(errorLogger kafka.Logger) Option {
	return func(k *kafka.Writer) {
		k.ErrorLogger = errorLogger
	}
}

// A transport used to send messages to kafka clusters.
//
// If nil, DefaultTransport is used.
func WithTransport(transport kafka.RoundTripper) Option {
	return func(k *kafka.Writer) {
		k.Transport = transport
	}
}

// AllowAutoTopicCreation notifies writer to create topic if missing.
func WithAllowAutoTopicCreation(allowAutoTopicCreation bool) Option {
	return func(k *kafka.Writer) {
		k.AllowAutoTopicCreation = allowAutoTopicCreation
	}
}
