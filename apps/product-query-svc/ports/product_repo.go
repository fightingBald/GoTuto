package ports

import (
    "context"

    "github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

type ProductRepo interface {
    GetByID(ctx context.Context, id int64) (*domain.Product, error)
    // Search returns the page of items and the total count matching the query
    Search(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error)
    Create(ctx context.Context, p *domain.Product) (int64, error)
}
