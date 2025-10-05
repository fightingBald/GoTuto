package outbound

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

// CommentRepository abstracts persistence for product comments.
type CommentRepository interface {
	CreateComment(ctx context.Context, comment *domain.Comment) (int64, error)
	GetCommentByID(ctx context.Context, id int64) (*domain.Comment, error)
	ListCommentsByProduct(ctx context.Context, productID int64) ([]domain.Comment, error)
	UpdateComment(ctx context.Context, comment *domain.Comment) error
	DeleteComment(ctx context.Context, id int64) error
}
