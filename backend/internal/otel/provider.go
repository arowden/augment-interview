package otel

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	ServiceName  string  `envconfig:"OTEL_SERVICE_NAME" default:"augment-fund-api"`
	Version      string  `envconfig:"VERSION" default:"dev"`
	Environment  string  `envconfig:"ENVIRONMENT" default:"development"`
	OTLPEndpoint string  `envconfig:"OTEL_EXPORTER_OTLP_ENDPOINT" default:""`
	SampleRate   float64 `envconfig:"OTEL_TRACES_SAMPLER_ARG" default:"0.1"`
	Enabled      bool    `envconfig:"OTEL_ENABLED" default:"false"`
}

type Providers struct {
	Tracer   trace.Tracer
	Meter    metric.Meter
	Logger   *slog.Logger
	Shutdown func(context.Context) error
}

func Init(ctx context.Context, cfg Config, log *slog.Logger) (*Providers, error) {
	if !cfg.Enabled || cfg.OTLPEndpoint == "" {
		log.Info("telemetry disabled, using no-op providers")
		return &Providers{
			Tracer:   otel.Tracer(cfg.ServiceName),
			Meter:    otel.Meter(cfg.ServiceName),
			Logger:   log,
			Shutdown: func(context.Context) error { return nil },
		}, nil
	}

	res, err := newResource(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	tracerProvider, traceShutdown, err := initTracer(ctx, cfg, res)
	if err != nil {
		return nil, fmt.Errorf("initializing tracer: %w", err)
	}

	meterProvider, metricsShutdown, err := initMeter(ctx, cfg, res)
	if err != nil {
		_ = traceShutdown(ctx)
		return nil, fmt.Errorf("initializing meter: %w", err)
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	shutdown := func(ctx context.Context) error {
		var errs []error
		if err := traceShutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if err := metricsShutdown(ctx); err != nil {
			errs = append(errs, err)
		}
		if len(errs) > 0 {
			return fmt.Errorf("shutdown errors: %v", errs)
		}
		return nil
	}

	log.Info("telemetry initialized",
		slog.String("service", cfg.ServiceName),
		slog.String("endpoint", cfg.OTLPEndpoint),
		slog.Float64("sample_rate", cfg.SampleRate),
	)

	return &Providers{
		Tracer:   tracerProvider.Tracer(cfg.ServiceName),
		Meter:    meterProvider.Meter(cfg.ServiceName),
		Logger:   log,
		Shutdown: shutdown,
	}, nil
}

func newResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	hostname, _ := os.Hostname()

	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.Version),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithAttributes(
			semconv.HostNameKey.String(hostname),
		),
	)
}

func initTracer(ctx context.Context, cfg Config, res *resource.Resource) (*sdktrace.TracerProvider, func(context.Context) error, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating trace exporter: %w", err)
	}

	sampler := newSampler(cfg.SampleRate)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	return tp, tp.Shutdown, nil
}

func initMeter(ctx context.Context, cfg Config, res *resource.Resource) (*sdkmetric.MeterProvider, func(context.Context) error, error) {
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating metric exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)

	otel.SetMeterProvider(mp)

	return mp, mp.Shutdown, nil
}
