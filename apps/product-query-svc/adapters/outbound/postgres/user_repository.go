package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGUserRepo struct{ pool *pgxpool.Pool }

var _ outbound.UserRepository = (*PGUserRepo)(nil)

func NewUserRepository(pool *pgxpool.Pool) outbound.UserRepository { return &PGUserRepo{pool: pool} }

func (r *PGUserRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	q, args, err := psql.Select("id", "name", "email", "created_at").From("users").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var u domain.User
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	u.CreatedAt = u.CreatedAt.UTC()
	return &u, nil
}
