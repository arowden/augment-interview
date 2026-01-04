## Context
Transfers are the core business operation moving units between owners. Each transfer must validate ownership, update two cap table entries, and record the transfer atomically. Transfer history provides audit capability.

## Goals / Non-Goals
- Goals: Atomic transfers, full validation, audit trail, idempotency protection, pessimistic locking, pagination
- Non-Goals: Batch transfers, scheduled transfers, transfer approval workflows

## Decisions
- Decision: Single Validate method combining all rules (format + ownership)
- Alternatives considered: Separate ValidateFormat/ValidateOwnership (requires two calls), validation in entity (couples entity to rules)
- Rationale: Single validation entry point reduces API surface and ensures all rules run together

- Decision: Transfer service modifies Ownership aggregate directly (cross-aggregate)
- Alternatives considered: Event-driven updates (eventual consistency unacceptable for financial data), handler-level coordination (duplicates locking logic)
- **Trade-off documented**: This violates DDD strict aggregate boundaries. Accepted because: single database, same bounded context, financial accuracy requires strong consistency. Alternative (saga/event-driven) rejected due to eventual consistency being unacceptable for fund unit accounting.

- Decision: Use database transaction with pessimistic locking (SELECT FOR UPDATE)
- Alternatives considered: Saga pattern (overkill for single-DB), optimistic locking (more conflicts under load)

- Decision: Idempotency key for duplicate request protection
- Alternatives considered: No idempotency (allows duplicate transfers on retry), application-level deduplication (race conditions)

- Decision: Functional options pattern for Service dependency injection
- Alternatives considered: Constructor injection only (less flexible for testing)

## Package Structure
```
internal/transfer/
  entity.go      - Transfer struct (immutable record)
  request.go     - TransferRequest value object with idempotency key
  validator.go   - Validator with combined validation method
  errors.go      - ErrInsufficientUnits, ErrSelfTransfer, ErrInvalidUnits, ErrDuplicateTransfer
  repository.go  - Repository interface with pagination
  postgres.go    - PostgresRepository implementation
  service.go     - Service with functional options DI
```

## Entity Design
```go
type Transfer struct {
    ID             uuid.UUID
    FundID         uuid.UUID
    FromOwner      string
    ToOwner        string
    Units          int
    IdempotencyKey *uuid.UUID
    TransferredAt  time.Time
}

type TransferRequest struct {
    FundID         uuid.UUID
    FromOwner      string
    ToOwner        string
    Units          int
    IdempotencyKey *uuid.UUID
}
```

## Validator Design
```go
type Validator struct{}

func (v *Validator) Validate(ctx context.Context, req TransferRequest, fromEntry *ownership.Entry) error {
    if strings.TrimSpace(req.FromOwner) == "" || strings.TrimSpace(req.ToOwner) == "" {
        return ErrInvalidOwner
    }
    if req.Units <= 0 {
        return ErrInvalidUnits
    }
    if strings.EqualFold(req.FromOwner, req.ToOwner) {
        return ErrSelfTransfer
    }
    if fromEntry == nil {
        return ErrOwnerNotFound
    }
    if fromEntry.Units < req.Units {
        return ErrInsufficientUnits
    }
    return nil
}
```

## Cross-Aggregate Transaction Design
**IMPORTANT: This documents an intentional DDD boundary crossing.**

Transfer execution necessarily modifies the Ownership aggregate. This is accepted because:
1. Single PostgreSQL database (no distributed transaction needed)
2. Same bounded context (fund management)
3. Financial accuracy requires ACID guarantees - eventual consistency is unacceptable
4. Units must always sum to fund total - cannot tolerate intermediate states

The alternative (event-driven with compensation) was rejected because:
- Transfer failures could leave inconsistent unit counts temporarily
- Compensation logic adds significant complexity
- Users expect immediate transfer confirmation

## Service Execution Flow with Explicit Locking SQL
```go
func (s *Service) ExecuteTransfer(ctx context.Context, req TransferRequest) (*Transfer, error) {
    if err := s.validator.ValidateBasic(req); err != nil {
        return nil, err
    }

    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return nil, fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    if req.IdempotencyKey != nil {
        existing, err := s.repo.FindByIdempotencyKey(ctx, tx, *req.IdempotencyKey)
        if err == nil && existing != nil {
            return existing, nil
        }
    }

    var fromEntry ownership.Entry
    err = tx.QueryRow(ctx, `
        SELECT id, fund_id, owner_name, units, acquired_at, updated_at
        FROM cap_table_entries
        WHERE fund_id = $1 AND owner_name = $2 AND deleted_at IS NULL
        FOR UPDATE
    `, req.FundID, req.FromOwner).Scan(
        &fromEntry.ID, &fromEntry.FundID, &fromEntry.OwnerName,
        &fromEntry.Units, &fromEntry.AcquiredAt, &fromEntry.UpdatedAt,
    )
    if err == pgx.ErrNoRows {
        return nil, ErrOwnerNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("lock from_owner: %w", err)
    }

    if fromEntry.Units < req.Units {
        return nil, ErrInsufficientUnits
    }

    _, err = tx.Exec(ctx, `
        UPDATE cap_table_entries
        SET units = units - $1, updated_at = NOW()
        WHERE id = $2
    `, req.Units, fromEntry.ID)
    if err != nil {
        return nil, fmt.Errorf("decrement from_owner: %w", err)
    }

    _, err = tx.Exec(ctx, `
        INSERT INTO cap_table_entries (fund_id, owner_name, units, acquired_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        ON CONFLICT (fund_id, owner_name) DO UPDATE
        SET units = cap_table_entries.units + $3, updated_at = NOW()
    `, req.FundID, req.ToOwner, req.Units)
    if err != nil {
        return nil, fmt.Errorf("upsert to_owner: %w", err)
    }

    transfer := &Transfer{
        ID:             uuid.New(),
        FundID:         req.FundID,
        FromOwner:      req.FromOwner,
        ToOwner:        req.ToOwner,
        Units:          req.Units,
        IdempotencyKey: req.IdempotencyKey,
        TransferredAt:  time.Now(),
    }

    if err := s.repo.CreateTx(ctx, tx, transfer); err != nil {
        return nil, fmt.Errorf("record transfer: %w", err)
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, fmt.Errorf("commit: %w", err)
    }

    return transfer, nil
}
```

## Service with Functional Options
```go
type Service struct {
    repo      Repository
    ownership ownership.Repository
    pool      *pgxpool.Pool
    validator *Validator
}

type ServiceOption func(*Service)

func WithRepository(r Repository) ServiceOption {
    return func(s *Service) { s.repo = r }
}

func WithOwnershipRepository(r ownership.Repository) ServiceOption {
    return func(s *Service) { s.ownership = r }
}

func WithPool(p *pgxpool.Pool) ServiceOption {
    return func(s *Service) { s.pool = p }
}

func NewService(opts ...ServiceOption) *Service {
    s := &Service{validator: &Validator{}}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## Repository Interface with Pagination
```go
type Repository interface {
    Create(ctx context.Context, transfer *Transfer) error
    CreateTx(ctx context.Context, tx pgx.Tx, transfer *Transfer) error
    FindByFundID(ctx context.Context, fundID uuid.UUID, limit, offset int) (*TransferList, error)
    FindByIdempotencyKey(ctx context.Context, tx pgx.Tx, key uuid.UUID) (*Transfer, error)
}

type TransferList struct {
    Transfers  []*Transfer
    TotalCount int
    Limit      int
    Offset     int
}
```

## Risks / Trade-offs
- SELECT FOR UPDATE may cause contention under load → Acceptable for expected volume, monitor lock wait times
- Cross-aggregate modification → Documented trade-off, accepted for ACID requirements
- No transfer reversal capability → Matches requirements, can add later as compensating transfer
- Idempotency key is optional → Clients should always provide for safety

## Open Questions
- None
