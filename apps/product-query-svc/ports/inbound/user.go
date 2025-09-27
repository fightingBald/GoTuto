package inbound

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// UserQueries exposes read-oriented user use cases for driving adapters.
type UserQueries interface {
	FetchByID(ctx context.Context, id int64) (*domain.User, error)
}
