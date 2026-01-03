package ownership

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service provides cap table domain operations.
type Service struct {
	repo Repository
}

// ServiceOption configures a Service.
type ServiceOption func(*Service)

// WithRepository sets the repository for the service.
func WithRepository(repo Repository) ServiceOption {
	return func(s *Service) {
		s.repo = repo
	}
}

// NewService creates a new Service with the provided options.
// Returns an error if no repository is configured.
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

// GetCapTable retrieves the cap table for a fund with pagination.
// Default limit is 100, maximum is 1000.
func (s *Service) GetCapTable(ctx context.Context, fundID uuid.UUID, params ListParams) (*CapTableView, error) {
	return s.repo.FindByFundID(ctx, fundID, params)
}

// GetOwnership retrieves a single owner's stake in a fund.
// Returns ErrOwnerNotFound (wrapped) if the owner is not in the cap table.
func (s *Service) GetOwnership(ctx context.Context, fundID uuid.UUID, ownerName string) (*Entry, error) {
	return s.repo.FindByFundAndOwner(ctx, fundID, ownerName)
}
