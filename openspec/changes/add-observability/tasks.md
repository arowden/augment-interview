## 1. Provider Setup with Resource Attributes
- [ ] 1.1 Create `internal/otel/provider.go` with provider initialization
- [ ] 1.2 Create Config struct (ServiceName, Version, Environment, OTLPEndpoint, SampleRate, Enabled)
- [ ] 1.3 Implement newResource with service.name, service.version, deployment.environment, host.name
- [ ] 1.4 Implement InitTracer with OTLP exporter
- [ ] 1.5 Implement InitMeter with OTLP exporter
- [ ] 1.6 Implement InitLogger with structured JSON output
- [ ] 1.7 Implement Shutdown for graceful cleanup

## 2. Sampling Strategy
- [ ] 2.1 Create `internal/otel/sampler.go` with custom sampler
- [ ] 2.2 Implement parent-based sampling (ParentBased + TraceIDRatioBased)
- [ ] 2.3 Default to 10% sampling rate
- [ ] 2.4 Always sample if parent is sampled
- [ ] 2.5 Support OTEL_TRACES_SAMPLER_ARG configuration

## 3. HTTP Middleware
- [ ] 3.1 Create `internal/otel/middleware.go` with HTTP middleware
- [ ] 3.2 Wrap router with otelhttp.NewHandler
- [ ] 3.3 Add span attributes for route and method
- [ ] 3.4 Implement healthCheckFilter to exclude /health, /ready
- [ ] 3.5 Set Error span status for 5xx responses

## 4. Metric Definitions with Units
- [ ] 4.1 Create `internal/otel/metrics.go` with application metrics
- [ ] 4.2 Define fund_created_total counter (unit: "1")
- [ ] 4.3 Define transfer_executed_total counter (unit: "1")
- [ ] 4.4 Define transfer_units_total counter (unit: "units")
- [ ] 4.5 Define http_request_duration_seconds histogram (unit: "s", buckets)
- [ ] 4.6 Define db_pool_size gauge (unit: "connections")
- [ ] 4.7 Define db_pool_active_connections gauge (unit: "connections")
- [ ] 4.8 Define db_pool_idle_connections gauge (unit: "connections")
- [ ] 4.9 Define db_pool_wait_count counter (unit: "1")
- [ ] 4.10 Define db_pool_wait_duration_seconds histogram (unit: "s", buckets)

## 5. Database Integration
- [ ] 5.1 Add otelpgx tracer to database pool configuration
- [ ] 5.2 Verify query spans appear as children of request spans
- [ ] 5.3 Implement periodic pool metrics collection goroutine
- [ ] 5.4 Call pool.Stat() to populate gauges

## 6. Configuration
- [ ] 6.1 Add OTEL_* environment variables to config
- [ ] 6.2 Support OTEL_EXPORTER_OTLP_ENDPOINT
- [ ] 6.3 Support OTEL_SERVICE_NAME
- [ ] 6.4 Support OTEL_TRACES_SAMPLER=parentbased_traceidratio
- [ ] 6.5 Support OTEL_TRACES_SAMPLER_ARG (default 0.1)
- [ ] 6.6 Support VERSION and ENVIRONMENT env vars
- [ ] 6.7 Support disabling telemetry for tests

## 7. SLO Documentation
- [ ] 7.1 Document p99 latency SLO (< 500ms)
- [ ] 7.2 Document p50 latency SLO (< 100ms)
- [ ] 7.3 Document transfer error rate SLO (< 1%)
- [ ] 7.4 Document pool wait duration SLO (p99 < 100ms)
- [ ] 7.5 Document pool saturation SLO (< 80%)

## 8. Local Development
- [ ] 8.1 Add Jaeger to docker-compose.yml
- [ ] 8.2 Configure OTLP endpoint for local Jaeger
- [ ] 8.3 Verify traces visible in Jaeger UI
- [ ] 8.4 Verify resource attributes appear in traces
