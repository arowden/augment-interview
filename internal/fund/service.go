package fund

import (
	"context"
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Service provides fund domain operations.
type Service struct {
	repo          Repository
	pool          *pgxpool.Pool
	ownershipRepo ownership.Repository
}

// ServiceOption configures a Service.
type ServiceOption func(*Service)

// WithPool sets the database connection pool for transactional operations.
func WithPool(p *pgxpool.Pool) ServiceOption {
	return func(s *Service) { s.pool = p }
}

// WithOwnershipRepository sets the ownership repository for fund creation.
func WithOwnershipRepository(r ownership.Repository) ServiceOption {
	return func(s *Service) { s.ownershipRepo = r }
}

// NewService creates a new Service with the required repository.
// Returns an error if repo is nil.
func NewService(repo Repository, opts ...ServiceOption) (*Service, error) {
	if repo == nil {
		return nil, errors.New("fund: repository is required")
	}
	s := &Service{repo: repo}
	for _, opt := range opts {
		opt(s)
	}
	return s, nil
}

// CreateFund creates a new fund with validation.
// Deprecated: Use CreateFundWithInitialOwner for transactional fund creation.
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

// CreateFundWithInitialOwner creates a new fund and assigns all units to the initial owner
// in a single atomic transaction. This ensures the fund and its initial cap table entry
// are always created together, maintaining data consistency.
//
// Returns ErrInvalidFund if fund validation fails.
// Returns ErrPoolRequired if pool or ownershipRepo is not configured.
func (s *Service) CreateFundWithInitialOwner(ctx context.Context, name string, totalUnits int, initialOwner string) (*Fund, error) {
	// Validate fund data first (fail fast before starting transaction).
	fund, err := NewFund(name, totalUnits)
	if err != nil {
		return nil, err
	}

	// Validate initial ownership entry.
	entry, err := ownership.NewCapTableEntry(fund.ID, initialOwner, totalUnits)
	if err != nil {
		return nil, fmt.Errorf("invalid initial owner: %w", err)
	}

	// Require pool and ownership repo for transactional operation.
	if s.pool == nil {
		return nil, ErrPoolRequired
	}
	if s.ownershipRepo == nil {
		return nil, ErrOwnershipRepoRequired
	}

	// Execute in a transaction.
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	// Create fund.
	if err := s.repo.CreateTx(ctx, tx, fund); err != nil {
		return nil, err
	}

	// Create initial ownership entry with all units.
	if err := s.ownershipRepo.CreateTx(ctx, tx, entry); err != nil {
		return nil, fmt.Errorf("create initial ownership: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return fund, nil
}

// GetFund retrieves a fund by ID.
func (s *Service) GetFund(ctx context.Context, id uuid.UUID) (*Fund, error) {
	return s.repo.FindByID(ctx, id)
}

// ListFunds retrieves funds with pagination.
func (s *Service) ListFunds(ctx context.Context, params ListParams) (*ListResult, error) {
	return s.repo.List(ctx, params)
}
