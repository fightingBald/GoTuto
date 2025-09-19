package app

import (
	"context"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

type ProductService struct {
	repo ports.ProductRepo
}

func NewProductService(r ports.ProductRepo) *ProductService { return &ProductService{repo: r} }

func (s *ProductService) GetProduct(ctx context.Context, id int64) (*domain.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) SearchProducts(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error) {
	return s.repo.Search(ctx, q, page, pageSize)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) CreateProduct(ctx context.Context, p *domain.Product) (int64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(ctx, p)
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	if p.ID <= 0 {
		return nil, domain.ErrValidation
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, p.ID)
}
