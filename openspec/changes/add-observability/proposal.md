# Change: Add Observability

## Why
Production systems require observability for debugging, performance monitoring, and operational insights. OpenTelemetry provides a vendor-neutral standard for traces, metrics, and logs with correlation across all signals.

## What Changes
- Add `internal/otel` package for OpenTelemetry setup
- Configure tracer, meter, and logger providers
- Add HTTP middleware for request tracing
- Configure OTLP exporters for traces and metrics
- Integrate with database layer for query tracing

## Impact
- Affected specs: observability (new)
- Affected code: `internal/otel/`, `internal/http/`, `internal/postgres/`
