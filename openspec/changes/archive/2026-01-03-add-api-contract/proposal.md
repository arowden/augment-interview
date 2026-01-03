# Change: Add OpenAPI Contract

## Why
Contract-first API development ensures frontend and backend teams work from a single source of truth. Generated code eliminates manual type synchronization and reduces integration bugs.

## What Changes
- Add OpenAPI 3.0 specification at `/api/openapi.yaml`
- Define all fund, cap-table, and transfer endpoints
- Define request/response schemas for all operations
- Enable code generation for Go server and TypeScript client

## Impact
- Affected specs: api-contract (new)
- Affected code: `/api/openapi.yaml`, generated server interfaces, generated client
