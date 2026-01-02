## Context
The Fund is the primary aggregate root in the cap table domain. It represents an investment vehicle with a fixed number of ownership units. Fund is a standalone aggregate - creating the initial ownership entry is a separate domain operation orchestrated at the handler level.

## Goals / Non-Goals
- Goals: Clean domain model, repository abstraction, aggregate isolation, constructor validation, dependency injection
- Non-Goals: Fund updates, fund deletion, multi-tenancy, coupling with ownership aggregate

## Decisions
- Decision: Fund and Ownership are separate aggregates with separate repositories
- Alternatives considered: Fund.Create with initialOwner (violates DDD single aggregate rule)
- Trade-off: Handler must coordinate both creates in a transaction, but aggregates remain decoupled

- Decision: Repository provides both Create and CreateTx methods for transaction support
- Alternatives considered: Only Create method (can't coordinate with other aggregates atomically)

- Decision: NewFund constructor returns (*Fund, error) with validation
- Alternatives considered: Validation only in service (domain entity should protect its invariants)

- Decision: Functional options pattern for Service dependency injection
- Alternatives considered: Constructor injection only (less flexible for testing)

- Decision: Repository interface in domain package, implementation in same package as postgres.go
- Alternatives considered: Separate infrastructure package (unnecessary indirection for this scope)

## Package Structure
```
internal/fund/
  entity.go      - Fund struct, NewFund() with validation
  errors.go      - ErrFundNotFound, ErrInvalidFund
  repository.go  - Repository interface with Tx support
  postgres.go    - PostgresRepository implementation
  service.go     - Service with functional options DI
```

## Entity Design
```go
type Fund struct {
    ID         uuid.UUID
    Name       string
    TotalUnits int
    CreatedAt  time.Time
}

func NewFund(name string, totalUnits int) (*Fund, error) {
    if strings.TrimSpace(name) == "" {
        return nil, ErrInvalidFund
    }
    if totalUnits <= 0 {
        return nil, ErrInvalidFund
    }
    return &Fund{
        ID:         uuid.New(),
        Name:       strings.TrimSpace(name),
        TotalUnits: totalUnits,
        CreatedAt:  time.Now(),
    }, nil
}
```

## Repository Interface
```go
type Repository interface {
    Create(ctx context.Context, fund *Fund) error
    CreateTx(ctx context.Context, tx pgx.Tx, fund *Fund) error
    FindByID(ctx context.Context, id uuid.UUID) (*Fund, error)
    FindAll(ctx context.Context) ([]*Fund, error)
}
```

## Service with Functional Options
```go
type Service struct {
    repo Repository
}

type ServiceOption func(*Service)

func WithRepository(r Repository) ServiceOption {
    return func(s *Service) {
        s.repo = r
    }
}

func NewService(opts ...ServiceOption) *Service {
    s := &Service{}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

func (s *Service) CreateFund(ctx context.Context, name string, totalUnits int) (*Fund, error) {
    fund, err := NewFund(name, totalUnits)
    if err != nil {
        return nil, err
    }
    if err := s.repo.Create(ctx, fund); err != nil {
        return nil, err
    }
    return fund, nil
}
```

## Handler-Level Orchestration
The HTTP handler coordinates fund and ownership creation in a single transaction:
```go
func (h *Handler) CreateFund(ctx context.Context, req CreateFundRequest) (*Fund, error) {
    tx, err := h.pool.Begin(ctx)
    if err != nil {
        return nil, err
    }
    defer tx.Rollback(ctx)

    fund, err := fund.NewFund(req.Name, req.TotalUnits)
    if err != nil {
        return nil, err
    }

    if err := h.fundRepo.CreateTx(ctx, tx, fund); err != nil {
        return nil, err
    }

    entry, err := ownership.NewCapTableEntry(fund.ID, req.InitialOwner, fund.TotalUnits)
    if err != nil {
        return nil, err
    }

    if err := h.ownershipRepo.CreateTx(ctx, tx, entry); err != nil {
        return nil, err
    }

    if err := tx.Commit(ctx); err != nil {
        return nil, err
    }

    return fund, nil
}
```

## Risks / Trade-offs
- Handler coordinates multiple aggregates → Acceptable, maintains aggregate boundaries
- No update/delete operations → Matches requirements, can add later if needed
- CreateTx exposes transaction to caller → Necessary for cross-aggregate coordination

## Open Questions
- None
