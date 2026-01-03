package fund

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// Service provides fund domain operations.
type Service struct {
	repo Repository
}

// ServiceOption configures a Service.
type ServiceOption func(*Service)

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

// GetFund retrieves a fund by ID.
func (s *Service) GetFund(ctx context.Context, id uuid.UUID) (*Fund, error) {
	return s.repo.FindByID(ctx, id)
}

// ListFunds retrieves funds with pagination.
func (s *Service) ListFunds(ctx context.Context, params ListParams) (*ListResult, error) {
	return s.repo.List(ctx, params)
}
