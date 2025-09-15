package service

import (
	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

type ProductService struct {
	repo ports.ProductRepo
}

func NewProductService(r ports.ProductRepo) *ProductService { return &ProductService{repo: r} }

func (s *ProductService) GetProduct(id int64) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) SearchProducts(q string, page, pageSize int) ([]domain.Product, error) {
	return s.repo.Search(q, page, pageSize)
}
