package fund

import (
	"context"
	"errors"
	"fmt"

	"github.com/arowden/augment-fund/internal/ownership"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo          Repository
	pool          *pgxpool.Pool
	ownershipRepo ownership.Repository
}

type ServiceOption func(*Service)

func WithPool(p *pgxpool.Pool) ServiceOption {
	return func(s *Service) { s.pool = p }
}

func WithOwnershipRepository(r ownership.Repository) ServiceOption {
	return func(s *Service) { s.ownershipRepo = r }
}

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

func (s *Service) CreateFundWithInitialOwner(ctx context.Context, name string, totalUnits int, initialOwner string) (*Fund, error) {
	fund, err := NewFund(name, totalUnits)
	if err != nil {
		return nil, err
	}

	entry, err := ownership.NewCapTableEntry(fund.ID, initialOwner, totalUnits)
	if err != nil {
		return nil, fmt.Errorf("invalid initial owner: %w", err)
	}

	if s.pool == nil {
		return nil, ErrPoolRequired
	}
	if s.ownershipRepo == nil {
		return nil, ErrOwnershipRepoRequired
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := s.repo.CreateTx(ctx, tx, fund); err != nil {
		return nil, err
	}

	if err := s.ownershipRepo.CreateTx(ctx, tx, entry); err != nil {
		return nil, fmt.Errorf("create initial ownership: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return fund, nil
}

func (s *Service) GetFund(ctx context.Context, id uuid.UUID) (*Fund, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ListFunds(ctx context.Context, params ListParams) (*ListResult, error) {
	return s.repo.List(ctx, params)
}
