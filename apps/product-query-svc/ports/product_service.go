package ports

import (
    "context"

    "github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

type ProductService interface {
    GetProduct(ctx context.Context, id int64) (*domain.Product, error)
    // SearchProducts returns items and total count
    SearchProducts(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error)
}
