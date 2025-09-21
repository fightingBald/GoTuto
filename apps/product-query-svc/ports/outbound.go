// outbound: infrastructure-facing interfaces (repositories, gateways)
package ports

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// ProductRepo is an outbound port abstracting persistence concerns.
type ProductRepo interface {
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	// Search returns the page of items and the total count matching the query
	Search(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error)
	Create(ctx context.Context, p *domain.Product) (int64, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, p *domain.Product) error
}

// UserRepo abstracts access to persistent user data.
type UserRepo interface {
	GetUserByID(ctx context.Context, id int64) (*domain.User, error)
}
