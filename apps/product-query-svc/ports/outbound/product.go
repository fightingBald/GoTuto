package outbound

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// ProductRepository abstracts persistence concerns for product aggregates.
type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	Search(ctx context.Context, query string, page, pageSize int) ([]domain.Product, int, error)
	Create(ctx context.Context, product *domain.Product) (int64, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
}
