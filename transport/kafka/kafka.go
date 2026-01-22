package kafka // import "go.microcore.dev/framework/transport/kafka"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
	"time"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
	"go.microcore.dev/framework/telemetry"
	"go.microcore.dev/framework/transport/kafka/reader"
	"go.microcore.dev/framework/transport/kafka/writer"

	"github.com/segmentio/kafka-go"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type (
	Manager interface {
		GetTopicWriter(topic string) *kafka.Writer
		GetTopicReader(topic string) *kafka.Reader
		GetTelemetryManager() telemetry.Manager
		SetTopicWriter(topic string, writer *kafka.Writer) Manager
		SetTopicReader(topic string, reader *kafka.Reader) Manager
		NewTopicWriter(topic string, opts ...writer.Option) Manager
		NewTopicReader(topic string, opts ...reader.Option) Manager
		SetTelemetryManager(telemetry telemetry.Manager) Manager
		Pub(topic string, payload []byte, opts ...PubOption) error
		PubJson(topic string, payload any, opts ...PubOption) error
		Sub(topic string, opts ...SubOption) Manager
		GetShutdownTimeout() time.Duration
		GetShutdownHandler() bool
		Shutdown(ctx context.Context, code int) error
	}

	Message                            = kafka.Message
	MessageHandler                     func(ctx context.Context, message kafka.Message) error
	PayloadParserMessageHandler[T any] func(ctx context.Context, message kafka.Message, payload *T) error

	k struct {
		brokers         *brokers
		writers         map[string]*kafka.Writer
		readers         map[string]*kafka.Reader
		telemetry       telemetry.Manager
		shutdownTimeout time.Duration
		shutdownHandler bool
		wg              sync.WaitGroup
	}

	brokers struct {
		writer []string
		reader []string
	}

	pub struct {
		context context.Context
		message kafka.Message
	}

	sub struct {
		context context.Context
		handler MessageHandler
	}

	headerCarrier struct {
		headers *[]kafka.Header
	}
)

var logger = log.New(pkg)

func New(opts ...Option) Manager {
	k := &k{
		brokers: &brokers{
			writer: []string{},
			reader: []string{},
		},
		writers:         make(map[string]*kafka.Writer),
		readers:         make(map[string]*kafka.Reader),
		shutdownTimeout: DefaultShutdownTimeout,
		shutdownHandler: DefaultShutdownHandler,
	}

	for _, opt := range opts {
		opt(k)
	}

	if k.shutdownHandler {
		shutdown.AddHandler(k.Shutdown)
		logger.Debug("shutdown handler has been successfully registered")
	}

	logger.Info(
		"created",
		slog.Group("shutdown",
			slog.Duration("timeout", k.shutdownTimeout),
			slog.Bool("handler", k.shutdownHandler),
		),
		slog.Bool("telemetry", k.telemetry != nil),
		slog.String("brokers.writer", strings.Join(k.brokers.writer, ",")),
		slog.String("brokers.reader", strings.Join(k.brokers.reader, ",")),
	)

	return k
}

func (k *k) GetTopicWriter(topic string) *kafka.Writer {
	return k.writers[topic]
}

func (k *k) GetTopicReader(topic string) *kafka.Reader {
	return k.readers[topic]
}

func (k *k) GetTelemetryManager() telemetry.Manager {
	return k.telemetry
}

func (k *k) SetTopicWriter(topic string, writer *kafka.Writer) Manager {
	k.writers[topic] = writer
	return k
}

func (k *k) SetTopicReader(topic string, reader *kafka.Reader) Manager {
	k.readers[topic] = reader
	return k
}

func (k *k) NewTopicWriter(topic string, opts ...writer.Option) Manager {
	defaults := []writer.Option{
		writer.WithAddr(
			kafka.TCP(k.brokers.writer...),
		),
		writer.WithTopic(topic),
	}
	k.writers[topic] = writer.New(
		append(defaults, opts...)...,
	)
	return k
}

func (k *k) NewTopicReader(topic string, opts ...reader.Option) Manager {
	defaults := []reader.Option{
		reader.WithBrokers(k.brokers.reader),
		reader.WithTopic(topic),
	}
	k.readers[topic] = reader.New(
		append(defaults, opts...)...,
	)
	return k
}

func (k *k) SetTelemetryManager(telemetry telemetry.Manager) Manager {
	k.telemetry = telemetry
	return k
}

func (k *k) Pub(topic string, payload []byte, opts ...PubOption) error {
	writer, ok := k.writers[topic]
	if !ok {
		return fmt.Errorf("writer for topic %q not found", topic)
	}
	pub := &pub{
		context: context.Background(),
		message: kafka.Message{
			Value:   payload,
			Headers: []kafka.Header{},
		},
	}
	for _, opt := range opts {
		opt(pub)
	}
	var span trace.Span
	if k.telemetry != nil {
		pub.context, span = k.telemetry.GetTracer().Start(pub.context, "kafka pub")
		defer span.End()
		k.telemetry.GetPropagator().Inject(pub.context, headerCarrier{&pub.message.Headers})
	}
	if err := writer.WriteMessages(pub.context, pub.message); err != nil {
		if k.telemetry != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	}
	return nil
}

func (k *k) PubJson(topic string, payload any, opts ...PubOption) error {
	p, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	return k.Pub(topic, p, opts...)
}

func (k *k) Sub(topic string, opts ...SubOption) Manager {
	reader, ok := k.readers[topic]
	if !ok {
		logger.Error(
			"sub: reader for topic not found",
			slog.String("topic", topic),
		)
		shutdown.Exit(shutdown.ExitUnavailable)
	}
	sub := &sub{
		context: context.Background(),
	}
	for _, opt := range opts {
		opt(sub)
	}
	if sub.handler == nil {
		logger.Error("sub: handler undefined")
		shutdown.Exit(shutdown.ExitSoftware)
	}
	k.wg.Add(1)
	logger.Info(
		"sub: reader started",
		slog.String("topic", topic),
	)
	go func() {
		defer k.wg.Done()
		for {
			msg, err := reader.ReadMessage(sub.context)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					logger.Info(
						"sub: context canceled, stop consuming topic",
						slog.String("topic", topic),
					)
					return
				}
				if errors.Is(err, io.EOF) {
					logger.Info(
						"sub: reader closed for topic",
						slog.String("topic", topic),
					)
					return
				}
				logger.Error(
					"sub: failed to read message from topic",
					slog.String("topic", topic),
					slog.Any("error", err),
				)
				continue
			}
			wctx := sub.context
			var span trace.Span
			if k.telemetry != nil {
				wctx = k.telemetry.GetPropagator().Extract(wctx, headerCarrier{&msg.Headers})
				wctx, span = k.telemetry.GetTracer().Start(wctx, "kafka sub")
			}
			if err := sub.handler(wctx, msg); err != nil {
				logger.Error(
					"sub: failed to handle message",
					slog.Any("error", err),
				)
				if k.telemetry != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
			}
			if k.telemetry != nil {
				span.End()
			}
		}
	}()
	return k
}

func (k *k) GetShutdownTimeout() time.Duration {
	return k.shutdownTimeout
}

func (k *k) GetShutdownHandler() bool {
	return k.shutdownHandler
}

func (k *k) Shutdown(ctx context.Context, code int) error {
	ctx, cancel := context.WithTimeout(ctx, k.shutdownTimeout)
	defer cancel()

	logger.Debug(
		"shutdown",
		slog.Int("code", code),
	)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	// Close writers
	for topic, w := range k.writers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := w.Close(); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("failed to close writer for topic %s: %w", topic, err))
				mu.Unlock()
			}
		}()
	}

	// Close readers
	for topic, r := range k.readers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.Close(); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("failed to close reader for topic %s: %w", topic, err))
				mu.Unlock()
			}
		}()
	}

	// Ждем внутренние горутины
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return errors.Join(errs...)
	}
}

func (f headerCarrier) Get(key string) string {
	for _, v := range *f.headers {
		if v.Key == key {
			return string(v.Value)
		}
	}
	return ""
}

func (f headerCarrier) Set(key, value string) {
	*f.headers = append(
		*f.headers,
		kafka.Header{
			Key:   key,
			Value: []byte(value),
		},
	)
}

func (f headerCarrier) Keys() []string {
	keys := []string{}
	for _, v := range *f.headers {
		keys = append(keys, string(v.Key))
	}
	return keys
}
