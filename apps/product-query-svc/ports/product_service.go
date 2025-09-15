package ports

import "github.com/fightingBald/GoTuto/apps/product-query-svc/domain"

type ProductService interface {
	GetProduct(id int64) (*domain.Product, error)
	SearchProducts(q string, page, pageSize int) ([]domain.Product, error)
}
