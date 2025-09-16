package service

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
