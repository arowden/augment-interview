package ownership

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

type ServiceOption func(*Service)

func WithRepository(repo Repository) ServiceOption {
	return func(s *Service) {
		s.repo = repo
	}
}

func NewService(opts ...ServiceOption) (*Service, error) {
	s := &Service{}
	for _, opt := range opts {
		opt(s)
	}
	if s.repo == nil {
		return nil, errors.New("ownership: repository is required")
	}
	return s, nil
}

func (s *Service) GetCapTable(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error) {
	return s.repo.FindByFundID(ctx, fundID, params)
}

func (s *Service) GetOwnership(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error) {
	return s.repo.FindByFundAndOwner(ctx, fundID, ownerName)
}

func (s *Service) CreateEntry(ctx context.Context, fundID uuid.UUID, ownerName string, units int) (*Entry, error) {
	entry, err := NewCapTableEntry(fundID, ownerName, units)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}
