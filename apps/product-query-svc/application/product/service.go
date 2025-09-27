package productapp

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/inbound"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
)

var _ inbound.ProductUseCases = (*Service)(nil)

// Service orchestrates product-related use cases across outbound dependencies.
type Service struct {
	repository outbound.ProductRepository
}

func NewService(repository outbound.ProductRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FetchByID(ctx context.Context, id int64) (*domain.Product, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) Search(ctx context.Context, query string, page, pageSize int) ([]domain.Product, int, error) {
	return s.repository.Search(ctx, query, page, pageSize)
}

func (s *Service) Remove(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *Service) Create(ctx context.Context, product *domain.Product) (int64, error) {
	if err := product.Validate(); err != nil {
		return 0, err
	}
	return s.repository.Create(ctx, product)
}

func (s *Service) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product.ID <= 0 {
		return nil, domain.ValidationError("id must be a positive integer")
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}
	if err := s.repository.Update(ctx, product); err != nil {
		return nil, err
	}
	return s.repository.GetByID(ctx, product.ID)
}
