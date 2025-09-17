// Package ports contains the stable boundaries of the core.
// inbound: application-facing use case interfaces
package ports

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// ProductService is an inbound port exposing application use cases
// to driving adapters (e.g., HTTP, gRPC, CLI).
type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*domain.Product, error)
	// SearchProducts returns items and total count
	SearchProducts(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error)
	// DeleteProduct removes the product by id
	DeleteProduct(ctx context.Context, id int64) error
}
