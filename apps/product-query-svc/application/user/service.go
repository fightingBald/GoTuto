package userapp

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/inbound"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
)

var _ inbound.UserQueries = (*Service)(nil)

// Service exposes user-specific use cases backed by a persistent repository.
type Service struct {
	repository outbound.UserRepository
}

func NewService(repository outbound.UserRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FetchByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ValidationError("id must be a positive integer")
	}
	return s.repository.FindByID(ctx, id)
}
