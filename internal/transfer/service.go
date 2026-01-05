package transfer

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
	ownershipRepo ownership.Repository
	pool          *pgxpool.Pool
	validator     *Validator
}

type ServiceOption func(*Service)

func WithRepository(r Repository) ServiceOption {
	return func(s *Service) { s.repo = r }
}

func WithOwnershipRepository(r ownership.Repository) ServiceOption {
	return func(s *Service) { s.ownershipRepo = r }
}

func WithPool(p *pgxpool.Pool) ServiceOption {
	return func(s *Service) { s.pool = p }
}

func NewService(opts ...ServiceOption) (*Service, error) {
	s := &Service{validator: NewValidator()}
	for _, opt := range opts {
		opt(s)
	}
	if s.repo == nil {
		return nil, errors.New("transfer: repository is required")
	}
	if s.ownershipRepo == nil {
		return nil, errors.New("transfer: ownership repository is required")
	}
	if s.pool == nil {
		return nil, errors.New("transfer: pool is required")
	}
	return s, nil
}

func (s *Service) ExecuteTransfer(ctx context.Context, req Request) (*Transfer, error) {
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
		if err != nil {
			return nil, fmt.Errorf("check idempotency key: %w", err)
		}
		if existing != nil {
			if existing.FundID != req.FundID ||
				existing.FromOwner != req.FromOwner ||
				existing.ToOwner != req.ToOwner ||
				existing.Units != req.Units {
				return nil, ErrDuplicateIdempotencyKey
			}
			return existing, nil
		}
	}

	fromEntry, err := s.ownershipRepo.FindByFundAndOwnerForUpdateTx(ctx, tx, req.FundID, req.FromOwner)
	if err != nil {
		if errors.Is(err, ownership.ErrOwnerNotFound) {
			return nil, ErrOwnerNotFound
		}
		return nil, fmt.Errorf("lock from_owner: %w", err)
	}

	if fromEntry.Units < req.Units {
		return nil, ErrInsufficientUnits
	}

	if err := s.ownershipRepo.DecrementUnitsTx(ctx, tx, fromEntry.ID, req.Units); err != nil {
		return nil, fmt.Errorf("decrement from_owner: %w", err)
	}

	if err := s.ownershipRepo.IncrementOrCreateTx(ctx, tx, req.FundID, req.ToOwner, req.Units); err != nil {
		return nil, fmt.Errorf("upsert to_owner: %w", err)
	}

	transfer := &Transfer{
		ID:             uuid.New(),
		FundID:         req.FundID,
		FromOwner:      req.FromOwner,
		ToOwner:        req.ToOwner,
		Units:          req.Units,
		IdempotencyKey: req.IdempotencyKey,
	}

	if err := s.repo.CreateTx(ctx, tx, transfer); err != nil {
		return nil, fmt.Errorf("record transfer: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return transfer, nil
}

func (s *Service) ListTransfers(ctx context.Context, fundID uuid.UUID, params ListParams) (*TransferList, error) {
	return s.repo.FindByFundID(ctx, fundID, params)
}
