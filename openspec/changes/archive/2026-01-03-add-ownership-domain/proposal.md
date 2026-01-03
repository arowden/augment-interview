# Change: Add Ownership Domain

## Why
The ownership domain represents the cap table - the authoritative record of who owns how many units in a fund. It provides read-only access to ownership state and supports transfer operations through controlled mutations.

## What Changes
- Add `internal/ownership` package with entry entity and repository
- Implement cap table retrieval by fund
- Support ownership adjustments for transfer operations
- Wire into HTTP handlers for cap table endpoint

## Impact
- Affected specs: ownership-domain (new)
- Affected code: `internal/ownership/`, `internal/http/handler.go`
