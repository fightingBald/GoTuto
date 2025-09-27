package inbound

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// ProductUseCases describes the application-facing entrypoints for product interactions.
type ProductUseCases interface {
	FetchByID(ctx context.Context, id int64) (*domain.Product, error)
	Search(ctx context.Context, query string, page, pageSize int) ([]domain.Product, int, error)
	Create(ctx context.Context, product *domain.Product) (int64, error)
	Update(ctx context.Context, product *domain.Product) (*domain.Product, error)
	Remove(ctx context.Context, id int64) error
}
