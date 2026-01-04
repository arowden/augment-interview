# Change: Add Transfer Domain

## Why
Transfers are the primary write operation in the cap table system. They move units between owners while maintaining invariants (no negative balances, no self-transfers, positive units). Transfer history provides an audit trail.

## What Changes
- Add `internal/transfer` package with transfer entity, validator, repository, and service
- Implement atomic transfer execution (update ownership, record transfer)
- Implement transfer history retrieval
- Wire into HTTP handlers for transfer endpoints

## Impact
- Affected specs: transfer-domain (new)
- Affected code: `internal/transfer/`, `internal/http/handler.go`
