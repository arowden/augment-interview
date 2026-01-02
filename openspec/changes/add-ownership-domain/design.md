## Context
Ownership entries form the cap table - each entry represents one owner's stake in a fund. The cap table is read-heavy but must support atomic updates during transfers. Ownership integrity is critical: the sum of all units must equal the fund's total units.

## Goals / Non-Goals
- Goals: Accurate ownership tracking, efficient cap table retrieval with pagination, atomic updates, transaction support for cross-aggregate coordination
- Non-Goals: Historical ownership tracking (handled by transfer history), percentage calculations (computed at presentation layer)

## Decisions
- Decision: Entry is the domain entity, CapTableView is a read model (not aggregate)
- Alternatives considered: CapTable as aggregate (misleading - it's a projection of entries, not a true aggregate with invariants)
- Rationale: CapTableView is read-only projection; Entry is the persisted entity modified by transfers

- Decision: Upsert operation for ownership changes (create if not exists, update if exists)
- Alternatives considered: Separate create/update (more code, same behavior)

- Decision: Repository provides both standalone and Tx methods for all write operations
- Alternatives considered: Only Tx methods (forces transaction even for simple cases)

- Decision: Pagination support via limit/offset on FindByFundID
- Alternatives considered: No pagination (fails for large cap tables), cursor pagination (more complex than needed)

- Decision: Functional options pattern for Service dependency injection
- Alternatives considered: Constructor injection only (less flexible for testing)

## Null Handling Behavior
- `FindByFundAndOwner` returns `ErrOwnerNotFound` when owner doesn't exist
- `FindByFundID` returns empty `CapTableView` with zero entries (never nil) for non-existent fund
- This distinction allows callers to differentiate "fund has no owners" from "owner not found"

## Package Structure
```
internal/ownership/
  entry.go       - Entry entity representing one ownership record
  captableview.go - CapTableView read model, TotalUnits() method
  errors.go      - ErrOwnerNotFound, ErrInsufficientUnits
  repository.go  - Repository interface with Tx support
  postgres.go    - PostgresRepository implementation
  service.go     - Service with functional options DI
```

## Entity Design
```go
type Entry struct {
    ID         uuid.UUID
    FundID     uuid.UUID
    OwnerName  string
    Units      int
    AcquiredAt time.Time
    UpdatedAt  time.Time
    DeletedAt  *time.Time
}

func NewCapTableEntry(fundID uuid.UUID, ownerName string, units int) (*Entry, error) {
    if strings.TrimSpace(ownerName) == "" {
        return nil, ErrInvalidOwner
    }
    if units < 0 {
        return nil, ErrInvalidUnits
    }
    return &Entry{
        ID:         uuid.New(),
        FundID:     fundID,
        OwnerName:  strings.TrimSpace(ownerName),
        Units:      units,
        AcquiredAt: time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}

type CapTableView struct {
    FundID     uuid.UUID
    Entries    []Entry
    TotalCount int
    Limit      int
    Offset     int
}

func (c *CapTableView) TotalUnits() int
func (c *CapTableView) FindOwner(name string) *Entry
```

## Repository Interface
```go
type Repository interface {
    Create(ctx context.Context, entry *Entry) error
    CreateTx(ctx context.Context, tx pgx.Tx, entry *Entry) error
    FindByFundID(ctx context.Context, fundID uuid.UUID, limit, offset int) (*CapTableView, error)
    FindByFundAndOwner(ctx context.Context, fundID uuid.UUID, owner string) (*Entry, error)
    Upsert(ctx context.Context, entry *Entry) error
    UpsertTx(ctx context.Context, tx pgx.Tx, entry *Entry) error
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

func (s *Service) GetCapTable(ctx context.Context, fundID uuid.UUID, limit, offset int) (*CapTableView, error) {
    if limit <= 0 {
        limit = 100
    }
    if limit > 1000 {
        limit = 1000
    }
    return s.repo.FindByFundID(ctx, fundID, limit, offset)
}
```

## Risks / Trade-offs
- Upsert requires careful handling of acquired_at vs updated_at → Only update updated_at on modifications
- Pagination adds complexity → Acceptable for API consistency and future scalability
- CapTableView naming → Clear that it's a read model, not a writable aggregate

## Open Questions
- None
