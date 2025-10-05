package inbound

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// CommentUseCases exposes comment workflows to inbound adapters.
type CommentUseCases interface {
	ListByProduct(ctx context.Context, productID int64) ([]domain.Comment, error)
	Create(ctx context.Context, productID, userID int64, content string) (*domain.Comment, error)
	Update(ctx context.Context, productID, commentID, userID int64, content string) (*domain.Comment, error)
	Delete(ctx context.Context, productID, commentID, userID int64) error
}
