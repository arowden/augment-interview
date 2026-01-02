## ADDED Requirements

### Requirement: OpenTelemetry Provider Initialization with Resource Attributes
The system SHALL initialize OpenTelemetry tracer, meter, and logger providers with service resource attributes.

#### Scenario: Provider initialization
- **WHEN** otel.Init is called with Config struct
- **THEN** Providers struct is returned with configured tracer, meter, and logger

#### Scenario: Resource attributes included
- **WHEN** providers are initialized
- **THEN** resource includes service.name, service.version, deployment.environment, and host.name

#### Scenario: OTLP endpoint configuration
- **WHEN** OTEL_EXPORTER_OTLP_ENDPOINT is set to "http://localhost:4317"
- **THEN** exporters send data to that endpoint

#### Scenario: Disabled telemetry
- **WHEN** OTEL_ENABLED is set to "false"
- **THEN** no-op providers are returned and no data is exported

### Requirement: Trace Sampling Strategy
The system SHALL use parent-based sampling with configurable sample rate.

#### Scenario: Default sample rate
- **WHEN** OTEL_TRACES_SAMPLER_ARG is not set
- **THEN** 10% of root spans are sampled

#### Scenario: Parent-based sampling
- **WHEN** incoming request has sampled traceparent header
- **THEN** the trace continues to be sampled (100%)

#### Scenario: Custom sample rate
- **WHEN** OTEL_TRACES_SAMPLER_ARG is set to 0.5
- **THEN** 50% of root spans are sampled

#### Scenario: Health check exclusion
- **WHEN** /health or /ready endpoints are called
- **THEN** no trace spans are created

### Requirement: Graceful Shutdown
The system SHALL flush and shutdown telemetry providers on application shutdown.

#### Scenario: Shutdown call
- **WHEN** Providers.Shutdown is called
- **THEN** all pending telemetry is flushed and exporters are closed

#### Scenario: Shutdown timeout
- **WHEN** shutdown context has a 5-second timeout
- **THEN** shutdown completes within that timeout

### Requirement: HTTP Request Tracing
The system SHALL automatically create spans for all HTTP requests.

#### Scenario: Request span creation
- **WHEN** an HTTP request is received
- **THEN** a span is created with http.method, http.route, and http.status_code attributes

#### Scenario: Span naming
- **WHEN** GET /api/funds is called
- **THEN** the span is named "GET /api/funds"

#### Scenario: Error span status
- **WHEN** an HTTP request returns 5xx
- **THEN** the span status is set to Error

#### Scenario: Trace context propagation
- **WHEN** a request includes traceparent header
- **THEN** the span is a child of the incoming trace

### Requirement: Database Query Tracing
The system SHALL create child spans for all database queries.

#### Scenario: Query span creation
- **WHEN** a SQL query is executed
- **THEN** a span is created as a child of the current request span

#### Scenario: Query span attributes
- **WHEN** a database span is created
- **THEN** it includes db.system, db.statement, and db.operation attributes

### Requirement: Metric Definitions with Units
The system SHALL define all metrics with explicit units and histogram bounds.

#### Scenario: HTTP request duration histogram
- **WHEN** http_request_duration_seconds is defined
- **THEN** it has unit "s" and bucket boundaries [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10]

#### Scenario: Fund creation counter
- **WHEN** fund_created_total is defined
- **THEN** it has unit "1" and description

#### Scenario: Transfer execution counter
- **WHEN** transfer_executed_total is defined
- **THEN** it has unit "1" and description

#### Scenario: Transfer units counter
- **WHEN** transfer_units_total is defined
- **THEN** it has unit "units" and description

### Requirement: Database Connection Pool Metrics
The system SHALL expose database connection pool metrics with appropriate types.

#### Scenario: Pool size gauge
- **WHEN** db_pool_size is defined
- **THEN** it is a gauge with unit "connections"

#### Scenario: Active connections gauge
- **WHEN** db_pool_active_connections is defined
- **THEN** it is a gauge with unit "connections"

#### Scenario: Idle connections gauge
- **WHEN** db_pool_idle_connections is defined
- **THEN** it is a gauge with unit "connections"

#### Scenario: Wait count counter
- **WHEN** db_pool_wait_count is defined
- **THEN** it is a counter with unit "1"

#### Scenario: Wait duration histogram
- **WHEN** db_pool_wait_duration_seconds is defined
- **THEN** it has unit "s" and bucket boundaries [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1]

### Requirement: Structured Logging with Trace Correlation
The system SHALL emit structured JSON logs with automatic trace ID injection.

#### Scenario: Log structure
- **WHEN** a log message is emitted
- **THEN** it is formatted as JSON with timestamp, level, message, and attributes

#### Scenario: Trace correlation
- **WHEN** logging within a traced request
- **THEN** the log includes trace_id and span_id fields

#### Scenario: Log levels
- **WHEN** logger is used
- **THEN** it supports Debug, Info, Warn, and Error levels

### Requirement: Middleware Integration
The system SHALL provide middleware that wraps the HTTP router.

#### Scenario: Middleware wrapping
- **WHEN** WrapHandler is called with a router
- **THEN** a traced handler is returned

#### Scenario: All API routes traced
- **WHEN** any /api/* endpoint is called
- **THEN** it produces a trace span

### Requirement: SLO Reference Documentation
The system SHALL document SLO targets for key metrics in design.md.

#### Scenario: Latency SLOs documented
- **WHEN** design.md is examined
- **THEN** it includes p99 and p50 latency targets

#### Scenario: Error rate SLOs documented
- **WHEN** design.md is examined
- **THEN** it includes transfer error rate target

#### Scenario: Pool saturation SLOs documented
- **WHEN** design.md is examined
- **THEN** it includes connection pool saturation target

### Requirement: Local Development Support
The system SHALL support local observability via Jaeger in docker-compose.

#### Scenario: Jaeger service in compose
- **WHEN** docker-compose up is run
- **THEN** Jaeger is available at http://localhost:16686

#### Scenario: Traces visible in Jaeger
- **WHEN** API requests are made locally
- **THEN** traces appear in the Jaeger UI
