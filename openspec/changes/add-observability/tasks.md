## 1. Provider Setup with Resource Attributes
- [x] 1.1 Create `internal/otel/provider.go` with provider initialization
- [x] 1.2 Create Config struct (ServiceName, Version, Environment, OTLPEndpoint, SampleRate, Enabled)
- [x] 1.3 Implement newResource with service.name, service.version, deployment.environment, host.name
- [x] 1.4 Implement InitTracer with OTLP exporter
- [x] 1.5 Implement InitMeter with OTLP exporter
- [x] 1.6 Implement InitLogger with structured JSON output
- [x] 1.7 Implement Shutdown for graceful cleanup

## 2. Sampling Strategy
- [x] 2.1 Create `internal/otel/sampler.go` with custom sampler
- [x] 2.2 Implement parent-based sampling (ParentBased + TraceIDRatioBased)
- [x] 2.3 Default to 10% sampling rate
- [x] 2.4 Always sample if parent is sampled
- [x] 2.5 Support OTEL_TRACES_SAMPLER_ARG configuration

## 3. HTTP Middleware
- [x] 3.1 Create `internal/otel/middleware.go` with HTTP middleware
- [x] 3.2 Wrap router with otelhttp.NewHandler
- [x] 3.3 Add span attributes for route and method
- [x] 3.4 Implement healthCheckFilter to exclude /health, /ready
- [x] 3.5 Set Error span status for 5xx responses

## 4. Metric Definitions with Units
- [x] 4.1 Create `internal/otel/metrics.go` with application metrics
- [x] 4.2 Define fund_created_total counter (unit: "1")
- [x] 4.3 Define transfer_executed_total counter (unit: "1")
- [x] 4.4 Define transfer_units_total counter (unit: "units")
- [x] 4.5 Define http_request_duration_seconds histogram (unit: "s", buckets)
- [x] 4.6 Define db_pool_size gauge (unit: "connections") - already in internal/postgres/metrics.go
- [x] 4.7 Define db_pool_active_connections gauge (unit: "connections") - already in internal/postgres/metrics.go
- [x] 4.8 Define db_pool_idle_connections gauge (unit: "connections") - already in internal/postgres/metrics.go
- [ ] 4.9 Define db_pool_wait_count counter (unit: "1") - deferred: requires pgx stat extensions
- [ ] 4.10 Define db_pool_wait_duration_seconds histogram (unit: "s", buckets) - deferred: requires pgx stat extensions

## 5. Database Integration
- [x] 5.1 Add otelpgx tracer to database pool configuration - already in internal/postgres/pool.go
- [x] 5.2 Verify query spans appear as children of request spans
- [x] 5.3 Implement periodic pool metrics collection goroutine - using observable gauges instead
- [x] 5.4 Call pool.Stat() to populate gauges - already in internal/postgres/metrics.go

## 6. Configuration
- [x] 6.1 Add OTEL_* environment variables to config
- [x] 6.2 Support OTEL_EXPORTER_OTLP_ENDPOINT
- [x] 6.3 Support OTEL_SERVICE_NAME
- [x] 6.4 Support OTEL_TRACES_SAMPLER=parentbased_traceidratio
- [x] 6.5 Support OTEL_TRACES_SAMPLER_ARG (default 0.1)
- [x] 6.6 Support VERSION and ENVIRONMENT env vars
- [x] 6.7 Support disabling telemetry for tests

## 7. SLO Documentation
- [x] 7.1 Document p99 latency SLO (< 500ms) - in design.md
- [x] 7.2 Document p50 latency SLO (< 100ms) - in design.md
- [x] 7.3 Document transfer error rate SLO (< 1%) - in design.md
- [x] 7.4 Document pool wait duration SLO (p99 < 100ms) - in design.md
- [x] 7.5 Document pool saturation SLO (< 80%) - in design.md

## 8. Local Development
- [x] 8.1 Add Jaeger to docker-compose.yml
- [x] 8.2 Configure OTLP endpoint for local Jaeger
- [ ] 8.3 Verify traces visible in Jaeger UI - manual verification required
- [ ] 8.4 Verify resource attributes appear in traces - manual verification required
