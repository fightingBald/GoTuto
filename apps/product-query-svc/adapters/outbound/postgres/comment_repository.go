package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGCommentRepo struct {
	pool *pgxpool.Pool
}

var _ outbound.CommentRepository = (*PGCommentRepo)(nil)

func NewCommentRepository(pool *pgxpool.Pool) outbound.CommentRepository {
	return &PGCommentRepo{pool: pool}
}

func (r *PGCommentRepo) CreateComment(ctx context.Context, comment *domain.Comment) (int64, error) {
	createdAt := comment.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	updatedAt := comment.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}

	qb := psql.Insert("comments").Columns(
		"product_id",
		"user_id",
		"content",
		"created_at",
		"updated_at",
	).Values(comment.ProductID, comment.UserID, comment.Content, createdAt, updatedAt).Suffix("RETURNING id")

	sql, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	if err := r.pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}

	comment.ID = id
	comment.CreatedAt = createdAt
	comment.UpdatedAt = updatedAt
	return id, nil
}

func (r *PGCommentRepo) GetCommentByID(ctx context.Context, id int64) (*domain.Comment, error) {
	qb := psql.Select("id", "product_id", "user_id", "content", "created_at", "updated_at").From("comments").Where(squirrel.Eq{"id": id})
	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	var c domain.Comment
	if err := r.pool.QueryRow(ctx, sql, args...).Scan(&c.ID, &c.ProductID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *PGCommentRepo) ListCommentsByProduct(ctx context.Context, productID int64) ([]domain.Comment, error) {
	qb := psql.Select("id", "product_id", "user_id", "content", "created_at", "updated_at").
		From("comments").
		Where(squirrel.Eq{"product_id": productID}).
		OrderBy("created_at DESC", "id DESC")

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PGCommentRepo) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	updatedAt := comment.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now().UTC()
	}

	qb := psql.Update("comments").
		Set("content", comment.Content).
		Set("updated_at", updatedAt).
		Where(squirrel.Eq{"id": comment.ID})

	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	ct, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	comment.UpdatedAt = updatedAt
	return nil
}

func (r *PGCommentRepo) DeleteComment(ctx context.Context, id int64) error {
	qb := psql.Delete("comments").Where(squirrel.Eq{"id": id})

	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	ct, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
