package outbound

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// UserRepository abstracts access to persistent user data.
type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*domain.User, error)
}
