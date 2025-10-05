package httpadapter

import (
	"context"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
)

func (s *Server) ListProductComments(ctx context.Context, request ListProductCommentsRequestObject) (ListProductCommentsResponseObject, error) {
	comments, err := s.comments.ListByProduct(ctx, request.ProductId)
	if err != nil {
		if resp, handled := listCommentsError(err); handled {
			return resp, nil
		}
		return nil, err
	}
	return okListComments(comments), nil
}

func (s *Server) CreateProductComment(ctx context.Context, request CreateProductCommentRequestObject) (CreateProductCommentResponseObject, error) {
	userID, content, err := commentCreateInput(request.Body)
	if err != nil {
		if resp, handled := createCommentError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	comment, err := s.comments.Create(ctx, request.ProductId, userID, content)
	if err != nil {
		if resp, handled := createCommentError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okCreateComment(comment), nil
}

func (s *Server) UpdateProductComment(ctx context.Context, request UpdateProductCommentRequestObject) (UpdateProductCommentResponseObject, error) {
	userID, content, err := commentUpdateInput(request.Body)
	if err != nil {
		if resp, handled := updateCommentError(err); handled {
			return resp, nil
		}
		return nil, err
	}
	if request.Params.UserId != 0 && request.Params.UserId != userID {
		if resp, handled := updateCommentError(domain.ValidationError("user id mismatch")); handled {
			return resp, nil
		}
		return nil, domain.ValidationError("user id mismatch")
	}

	updated, err := s.comments.Update(ctx, request.ProductId, request.CommentId, userID, content)
	if err != nil {
		if resp, handled := updateCommentError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okUpdateComment(updated), nil
}

func (s *Server) DeleteProductComment(ctx context.Context, request DeleteProductCommentRequestObject) (DeleteProductCommentResponseObject, error) {
	if err := s.comments.Delete(ctx, request.ProductId, request.CommentId, request.Params.UserId); err != nil {
		if resp, handled := deleteCommentError(err); handled {
			return resp, nil
		}
		return nil, err
	}

	return okDeleteComment(), nil
}
