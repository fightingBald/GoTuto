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
	// CreateProduct validates and persists a new product, returning its id
	CreateProduct(ctx context.Context, p *domain.Product) (int64, error)
	// UpdateProduct replaces the existing product state and returns the updated snapshot
	UpdateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error)
}

// UserService exposes user-related use cases to driving adapters.
type UserService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
}
