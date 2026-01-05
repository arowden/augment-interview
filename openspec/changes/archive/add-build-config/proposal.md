# Change: Add Build Configuration

## Why
Using Datadog's Orchestrion toolchain enables automatic compile-time instrumentation for OpenTelemetry. This eliminates manual instrumentation boilerplate while providing comprehensive tracing and logging across the entire application without code changes.

## What Changes
- Add Dockerfile with multi-stage build using Orchestrion
- Configure go build with `-toolexec` for compile-time instrumentation
- Add Makefile targets for instrumented and non-instrumented builds
- Configure automatic HTTP, database, and runtime instrumentation

## Impact
- Affected specs: build-config (new)
- Affected code: `/Dockerfile`, `/Makefile`, `/.orchestrion.yml`
