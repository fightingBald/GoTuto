package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type PGProductRepo struct{ pool *pgxpool.Pool }

var _ outbound.ProductRepository = (*PGProductRepo)(nil)

func NewProductRepository(pool *pgxpool.Pool) outbound.ProductRepository {
	return &PGProductRepo{pool: pool}
}

func (r *PGProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	q, args, err := psql.Select("id", "name", "price", "tags").From("products").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	var p domain.Product
	var tags []string
	if err := r.pool.QueryRow(ctx, q, args...).Scan(&p.ID, &p.Name, &p.Price, &tags); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

func (r *PGProductRepo) Search(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	base := psql.Select("id", "name", "price", "tags").From("products")
	countB := psql.Select("COUNT(*)").From("products")
	if strings.TrimSpace(q) != "" {
		base = base.Where("name ILIKE ?", "%"+q+"%")
		countB = countB.Where("name ILIKE ?", "%"+q+"%")
	}
	builder := base.OrderBy("id").Limit(uint64(pageSize)).Offset(uint64(offset))

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.Product
	for rows.Next() {
		var p domain.Product
		var tags []string
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &tags); err != nil {
			return nil, 0, err
		}
		p.Tags = tags
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// total count
	cq, cargs, err := countB.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.pool.QueryRow(ctx, cq, cargs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *PGProductRepo) Create(ctx context.Context, p *domain.Product) (int64, error) {
	ib := psql.Insert("products").Columns("name", "price", "tags").Values(p.Name, p.Price, p.Tags).Suffix("RETURNING id")
	sql, args, err := ib.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	if err := r.pool.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PGProductRepo) Delete(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, "DELETE FROM products WHERE id=$1", id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PGProductRepo) Update(ctx context.Context, p *domain.Product) error {
	ct, err := r.pool.Exec(ctx, "UPDATE products SET name=$1, price=$2, tags=$3 WHERE id=$4", p.Name, p.Price, p.Tags, p.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
