package telemetry // import "go.microcore.dev/framework/telemetry"

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	logSdk "go.opentelemetry.io/otel/sdk/log"
	metricSdk "go.opentelemetry.io/otel/sdk/metric"
	metricSdkExemplar "go.opentelemetry.io/otel/sdk/metric/exemplar"
	traceSdk "go.opentelemetry.io/otel/sdk/trace"

	telemetryLog "go.microcore.dev/framework/telemetry/log"

	logProvider "go.microcore.dev/framework/telemetry/log/provider"
	metricProvider "go.microcore.dev/framework/telemetry/metric/provider"
	traceProvider "go.microcore.dev/framework/telemetry/trace/provider"

	logOtlpGrpcExporter "go.microcore.dev/framework/telemetry/log/exporter/otlp/grpc"
	metricOtlpGrpcExporter "go.microcore.dev/framework/telemetry/metric/exporter/otlp/grpc"
	traceOtlpGrpcExporter "go.microcore.dev/framework/telemetry/trace/exporter/otlp/grpc"

	metricPeriodicReader "go.microcore.dev/framework/telemetry/metric/reader/periodic"

	otelLog "go.opentelemetry.io/otel/log"
	otelMetric "go.opentelemetry.io/otel/metric"
	otelTrace "go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	_ "go.microcore.dev/framework"
	"go.microcore.dev/framework/log"
	"go.microcore.dev/framework/shutdown"
)

type (
	Manager interface {
		GetTraceProvider() *traceSdk.TracerProvider
		GetMetricProvider() *metricSdk.MeterProvider
		GetLogProvider() *logSdk.LoggerProvider
		GetTracer() otelTrace.Tracer
		GetMeter() otelMetric.Meter
		GetLogger() otelLog.Logger
		GetPropagator() propagation.TextMapPropagator
		GetMetricsHttpHandler() http.Handler
		GetShutdownTimeout() time.Duration
		GetShutdownHandler() bool
		ForceFlush(ctx context.Context, reason string) error
		Shutdown(ctx context.Context, reason string) error
	}

	t struct {
		traceProvider   *traceSdk.TracerProvider
		metricProvider  *metricSdk.MeterProvider
		logProvider     *logSdk.LoggerProvider
		propagator      propagation.TextMapPropagator
		shutdownTimeout time.Duration
		shutdownHandler bool
		setLogProvider  bool
	}
)

var logger = log.New(pkg)

func New(opts ...Option) Manager {
	t := &t{
		propagator:      defaultPropagator,
		shutdownTimeout: DefaultShutdownTimeout,
		shutdownHandler: DefaultShutdownHandler,
		setLogProvider:  DefaultSetLogProvider,
	}

	for _, opt := range opts {
		opt(t)
	}

	if t.traceProvider == nil {
		t.traceProvider = traceProvider.New()
	}

	if t.metricProvider == nil {
		t.metricProvider = metricProvider.New()
	}

	if t.logProvider == nil {
		t.logProvider = logProvider.New()
	}

	if t.shutdownHandler {
		shutdown.AddHandler(t.Shutdown)
		logger.Debug("shutdown handler has been successfully registered")
	}

	if t.setLogProvider {
		slog.SetDefault(
			otelslog.NewLogger(
				logProvider.InstrumentationName,
				otelslog.WithLoggerProvider(t.logProvider),
			),
		)
		logger.Info(
			"logging backend has been successfully changed to otel",
			slog.String("instrumentation_name", logProvider.InstrumentationName),
		)
	}

	logger.Info(
		"manager has been successfully created",
		slog.Group("shutdown",
			slog.Duration("timeout", t.shutdownTimeout),
			slog.Bool("handler", t.shutdownHandler),
		),
		slog.Bool("default_log_provider", t.setLogProvider),
	)

	return t
}

func NewDefaultInsecureOtlpGrpc(ctx context.Context, endpoint string, service string) Manager {
	host, err := os.Hostname()
	if err != nil {
		host = "undefined"
	}

	return New(
		WithTraceProviderOptions(
			traceProvider.WithBatcher(
				traceOtlpGrpcExporter.New(
					ctx,
					traceOtlpGrpcExporter.WithEndpoint(endpoint),
					traceOtlpGrpcExporter.WithInsecure(),
				),
			),
			traceProvider.WithSampler(
				traceSdk.ParentBased(
					traceSdk.AlwaysSample(),
					traceSdk.WithRemoteParentSampled(
						traceSdk.AlwaysSample(),
					),
					traceSdk.WithRemoteParentNotSampled(
						traceSdk.NeverSample(),
					),
					traceSdk.WithLocalParentSampled(
						traceSdk.AlwaysSample(),
					),
					traceSdk.WithLocalParentNotSampled(
						traceSdk.NeverSample(),
					),
				),
			),
			traceProvider.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(service),
					semconv.HostNameKey.String(host),
				),
			),
		),
		WithMetricProviderOptions(
			metricProvider.WithReader(
				metricPeriodicReader.New(
					metricOtlpGrpcExporter.New(
						ctx,
						metricOtlpGrpcExporter.WithEndpoint(endpoint),
						metricOtlpGrpcExporter.WithInsecure(),
					),
					metricPeriodicReader.WithInterval(DefaultMetricPeriodicReaderInterval),
					metricPeriodicReader.WithTimeout(DefaultMetricPeriodicReaderTimeout),
				),
			),
			metricProvider.WithExemplarFilter(
				metricSdkExemplar.TraceBasedFilter,
			),
			metricProvider.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(service),
					semconv.HostNameKey.String(host),
				),
			),
		),
		WithLogProviderOptions(
			logProvider.WithProcessor(
				telemetryLog.NewProcessor(
					logSdk.NewBatchProcessor(
						logOtlpGrpcExporter.New(
							ctx,
							logOtlpGrpcExporter.WithEndpoint(endpoint),
							logOtlpGrpcExporter.WithInsecure(),
						),
						logSdk.WithExportInterval(DefaultLogExportInterval),
						logSdk.WithExportTimeout(DefaultLogExportTimeout),
					),
				),
			),
			logProvider.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(service),
					semconv.HostNameKey.String(host),
				),
			),
		),
	)
}

func (t *t) GetTraceProvider() *traceSdk.TracerProvider {
	return t.traceProvider
}

func (t *t) GetMetricProvider() *metricSdk.MeterProvider {
	return t.metricProvider
}

func (t *t) GetLogProvider() *logSdk.LoggerProvider {
	return t.logProvider
}

func (t *t) GetTracer() otelTrace.Tracer {
	return t.traceProvider.Tracer(traceProvider.InstrumentationName)
}

func (t *t) GetMeter() otelMetric.Meter {
	return t.metricProvider.Meter(metricProvider.InstrumentationName)
}

func (t *t) GetLogger() otelLog.Logger {
	return t.logProvider.Logger(logProvider.InstrumentationName)
}

func (t *t) GetPropagator() propagation.TextMapPropagator {
	return t.propagator
}

func (t *t) GetMetricsHttpHandler() http.Handler {
	return promhttp.Handler()
}

func (t *t) GetShutdownTimeout() time.Duration {
	return t.shutdownTimeout
}

func (t *t) GetShutdownHandler() bool {
	return t.shutdownHandler
}

func (t *t) ForceFlush(ctx context.Context, reason string) error {
	logger.Info(
		"force flush",
		slog.String("reason", reason),
	)

	if t.traceProvider != nil {
		if err := t.traceProvider.ForceFlush(ctx); err != nil {
			return fmt.Errorf("telemetry trace force flush failed: %v", err)
		}
	}
	if t.metricProvider != nil {
		if err := t.metricProvider.ForceFlush(ctx); err != nil {
			return fmt.Errorf("telemetry metric force flush failed: %v", err)
		}
	}
	if t.logProvider != nil {
		if err := t.logProvider.ForceFlush(ctx); err != nil {
			return fmt.Errorf("telemetry log force flush failed: %v", err)
		}
	}

	return nil
}

func (t *t) Shutdown(ctx context.Context, reason string) error {
	ctx, cancel := context.WithTimeout(ctx, t.shutdownTimeout)
	defer cancel()

	logger.Debug(
		"shutdown",
		slog.String("reason", reason),
	)

	if err := t.ForceFlush(ctx, reason); err != nil {
		return err
	}

	if t.traceProvider != nil {
		if err := t.traceProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("telemetry trace shutdown failed: %v", err)
		}
	}
	if t.metricProvider != nil {
		if err := t.metricProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("telemetry metric shutdown failed: %v", err)
		}
	}
	if t.logProvider != nil {
		if err := t.logProvider.Shutdown(ctx); err != nil {
			return fmt.Errorf("telemetry log shutdown failed: %v", err)
		}
	}

	return nil
}


