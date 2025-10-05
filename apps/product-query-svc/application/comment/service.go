package commentapp

import (
	"context"
	"strings"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/inbound"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
)

var _ inbound.CommentUseCases = (*Service)(nil)

// Service coordinates comment operations against domain rules and persistence.
type Service struct {
	comments outbound.CommentRepository
	products outbound.ProductRepository
	users    outbound.UserRepository
}

func NewService(comments outbound.CommentRepository, products outbound.ProductRepository, users outbound.UserRepository) *Service {
	return &Service{comments: comments, products: products, users: users}
}

func (s *Service) ListByProduct(ctx context.Context, productID int64) ([]domain.Comment, error) {
	if productID <= 0 {
		return nil, domain.ValidationError("product id must be a positive integer")
	}
	if _, err := s.products.GetByID(ctx, productID); err != nil {
		return nil, err
	}
	return s.comments.ListCommentsByProduct(ctx, productID)
}

func (s *Service) Create(ctx context.Context, productID, userID int64, content string) (*domain.Comment, error) {
	if productID <= 0 {
		return nil, domain.ValidationError("product id must be a positive integer")
	}
	if userID <= 0 {
		return nil, domain.ValidationError("user id must be a positive integer")
	}

	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, domain.ValidationError("content required")
	}

	if _, err := s.products.GetByID(ctx, productID); err != nil {
		return nil, err
	}
	if _, err := s.users.FindByID(ctx, userID); err != nil {
		return nil, err
	}

	comment, err := domain.NewComment(productID, userID, trimmed)
	if err != nil {
		return nil, err
	}

	id, err := s.comments.CreateComment(ctx, comment)
	if err != nil {
		return nil, err
	}
	comment.ID = id
	return comment, nil
}

func (s *Service) Update(ctx context.Context, productID, commentID, userID int64, content string) (*domain.Comment, error) {
	if productID <= 0 {
		return nil, domain.ValidationError("product id must be a positive integer")
	}
	if commentID <= 0 {
		return nil, domain.ValidationError("comment id must be a positive integer")
	}
	if userID <= 0 {
		return nil, domain.ValidationError("user id must be a positive integer")
	}

	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil, domain.ValidationError("content required")
	}

	existing, err := s.comments.GetCommentByID(ctx, commentID)
	if err != nil {
		return nil, err
	}
	if existing.ProductID != productID {
		return nil, domain.ErrNotFound
	}
	if existing.UserID != userID {
		return nil, domain.ForbiddenError("cannot modify another user's comment")
	}

	if err := existing.UpdateContent(trimmed); err != nil {
		return nil, err
	}
	if err := s.comments.UpdateComment(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, productID, commentID, userID int64) error {
	if productID <= 0 {
		return domain.ValidationError("product id must be a positive integer")
	}
	if commentID <= 0 {
		return domain.ValidationError("comment id must be a positive integer")
	}
	if userID <= 0 {
		return domain.ValidationError("user id must be a positive integer")
	}

	existing, err := s.comments.GetCommentByID(ctx, commentID)
	if err != nil {
		return err
	}
	if existing.ProductID != productID {
		return domain.ErrNotFound
	}
	if existing.UserID != userID {
		return domain.ForbiddenError("cannot delete another user's comment")
	}

	return s.comments.DeleteComment(ctx, commentID)
}
