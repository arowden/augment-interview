## Context
OpenTelemetry is the industry standard for observability. It provides unified APIs for traces, metrics, and logs with automatic context propagation. OTLP exporters send data to various backends (Jaeger locally, AWS X-Ray in production).

## Goals / Non-Goals
- Goals: Distributed tracing with sampling, request metrics, database query visibility, structured logging with correlation, connection pool observability
- Non-Goals: Custom dashboards, alerting rules, log aggregation backend setup

## Decisions
- Decision: Use otelhttp middleware for automatic HTTP span creation
- Alternatives considered: Manual span creation (error-prone, verbose)

- Decision: Use otelpgx for automatic database query tracing
- Alternatives considered: Manual query wrapping (verbose), no DB tracing (reduced visibility)

- Decision: OTLP/gRPC exporters for both traces and metrics
- Alternatives considered: Jaeger-native exporter (vendor lock-in), Prometheus scraping (pull model complexity)

- Decision: Structured JSON logging with trace correlation via slog
- Alternatives considered: Zerolog (another dependency), logrus (deprecated), plain text (no structure)

- Decision: Parent-based sampling with 10% default, 100% for errors
- Alternatives considered: Always-on (too expensive at scale), head-based only (misses errors)

## Package Structure
```
internal/otel/
  provider.go    - Tracer, meter, logger initialization with resource
  middleware.go  - HTTP middleware wrapper
  metrics.go     - Application metric definitions with units
  sampler.go     - Custom sampler configuration
```

## Resource Attributes
```go
func newResource(ctx context.Context) (*resource.Resource, error) {
    return resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceNameKey.String(os.Getenv("OTEL_SERVICE_NAME")),
            semconv.ServiceVersionKey.String(os.Getenv("VERSION")),
            semconv.DeploymentEnvironmentKey.String(os.Getenv("ENVIRONMENT")),
            attribute.String("host.name", hostname),
        ),
        resource.WithOS(),
        resource.WithHost(),
    )
}
```

## Provider Initialization
```go
type Providers struct {
    Tracer   trace.Tracer
    Meter    metric.Meter
    Logger   *slog.Logger
    Shutdown func(context.Context) error
}

func Init(ctx context.Context, cfg Config) (*Providers, error)

type Config struct {
    ServiceName  string
    Version      string
    Environment  string
    OTLPEndpoint string
    SampleRate   float64
    Enabled      bool
}
```

## Sampling Strategy
```go
func newSampler(sampleRate float64) sdktrace.Sampler {
    return sdktrace.ParentBased(
        sdktrace.TraceIDRatioBased(sampleRate),
        sdktrace.WithRemoteParentSampled(sdktrace.AlwaysSample()),
        sdktrace.WithRemoteParentNotSampled(sdktrace.TraceIDRatioBased(sampleRate)),
    )
}
```

Environment configuration:
```
OTEL_TRACES_SAMPLER=parentbased_traceidratio
OTEL_TRACES_SAMPLER_ARG=0.1  # 10% default sampling
```

## HTTP Middleware
```go
func WrapHandler(handler http.Handler, operation string) http.Handler {
    return otelhttp.NewHandler(handler, operation,
        otelhttp.WithSpanNameFormatter(spanNameFormatter),
        otelhttp.WithFilter(healthCheckFilter),
    )
}

func healthCheckFilter(r *http.Request) bool {
    return r.URL.Path != "/health" && r.URL.Path != "/ready"
}
```

## Metric Definitions with Units
```go
var (
    FundCreatedTotal metric.Int64Counter
    TransferExecutedTotal metric.Int64Counter
    TransferUnitsTotal metric.Int64Counter
    HTTPRequestDuration metric.Float64Histogram
    DBPoolSize metric.Int64Gauge
    DBPoolActiveConnections metric.Int64Gauge
    DBPoolIdleConnections metric.Int64Gauge
    DBPoolWaitCount metric.Int64Counter
    DBPoolWaitDuration metric.Float64Histogram
)

func InitMetrics(meter metric.Meter) error {
    var err error

    FundCreatedTotal, err = meter.Int64Counter(
        "fund_created_total",
        metric.WithDescription("Total number of funds created"),
        metric.WithUnit("1"),
    )
    if err != nil {
        return err
    }

    TransferExecutedTotal, err = meter.Int64Counter(
        "transfer_executed_total",
        metric.WithDescription("Total number of transfers executed"),
        metric.WithUnit("1"),
    )
    if err != nil {
        return err
    }

    TransferUnitsTotal, err = meter.Int64Counter(
        "transfer_units_total",
        metric.WithDescription("Total units transferred"),
        metric.WithUnit("units"),
    )
    if err != nil {
        return err
    }

    HTTPRequestDuration, err = meter.Float64Histogram(
        "http_request_duration_seconds",
        metric.WithDescription("HTTP request latency in seconds"),
        metric.WithUnit("s"),
        metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
    )
    if err != nil {
        return err
    }

    DBPoolSize, err = meter.Int64Gauge(
        "db_pool_size",
        metric.WithDescription("Maximum size of connection pool"),
        metric.WithUnit("connections"),
    )
    if err != nil {
        return err
    }

    DBPoolActiveConnections, err = meter.Int64Gauge(
        "db_pool_active_connections",
        metric.WithDescription("Number of currently active connections"),
        metric.WithUnit("connections"),
    )
    if err != nil {
        return err
    }

    DBPoolWaitDuration, err = meter.Float64Histogram(
        "db_pool_wait_duration_seconds",
        metric.WithDescription("Time spent waiting for a connection"),
        metric.WithUnit("s"),
        metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1),
    )
    if err != nil {
        return err
    }

    return nil
}
```

## SLO Definitions (Monitoring Reference)
These SLOs inform alerting configuration (outside this implementation scope):

| Metric | SLO | Description |
|--------|-----|-------------|
| http_request_duration_seconds p99 | < 500ms | API latency |
| http_request_duration_seconds p50 | < 100ms | API latency |
| transfer_executed_total (error rate) | < 1% | Transfer success rate |
| db_pool_wait_duration_seconds p99 | < 100ms | DB connection acquisition |
| db_pool_active_connections / db_pool_size | < 80% | Pool saturation |

## Environment Configuration
```
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=augment-fund-api
OTEL_TRACES_SAMPLER=parentbased_traceidratio
OTEL_TRACES_SAMPLER_ARG=0.1
OTEL_ENABLED=true
VERSION=dev
ENVIRONMENT=development
```

## Trace Context Flow
1. HTTP request arrives
2. otelhttp extracts/creates trace context
3. Sampler decides whether to sample (10% default, 100% if parent sampled)
4. Context propagates through service layer
5. pgx spans attach as children to request span
6. Response completes, span ends with status

## Risks / Trade-offs
- OTLP requires running collector/backend → Mitigated by optional disabling
- 10% sampling may miss important traces → Use trace_id in logs for correlation
- Histogram buckets may need tuning → Start with reasonable defaults, adjust based on p99

## Open Questions
- None
