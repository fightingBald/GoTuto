package db

import (
	"context"
	"github.com/fightingBald/GoTuto/internal/domain"
	"github.com/fightingBald/GoTuto/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type PGProductRepo struct{ pool *pgxpool.Pool }

func NewPGProductRepo(pool *pgxpool.Pool) ports.ProductRepo { return &PGProductRepo{pool: pool} }

func (r *PGProductRepo) GetByID(id int64) (*domain.Product, error) {
	const q = `SELECT id,name,price,tags FROM products WHERE id=$1`
	var p domain.Product
	var tags []string
	if err := r.pool.QueryRow(context.Background(), q, id).Scan(&p.ID, &p.Name, &p.Price, &tags); err != nil {
		return nil, err
	}
	p.Tags = tags
	return &p, nil
}

func (r *PGProductRepo) Search(q string, page, pageSize int) ([]domain.Product, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize
	var args []any
	where := ``
	if strings.TrimSpace(q) != "" {
		where = `WHERE name ILIKE $1`
		args = append(args, "%"+q+"%")
	}
	sql := `SELECT id,name,price,tags FROM products ` + where + ` ORDER BY id LIMIT $2 OFFSET $3`
	args = append(args, pageSize, offset)

	rows, err := r.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Product
	for rows.Next() {
		var p domain.Product
		var tags []string
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &tags); err != nil {
			return nil, err
		}
		p.Tags = tags
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PGProductRepo) Create(p *domain.Product) (int64, error) {
	const q = `INSERT INTO products(name, price, tags) VALUES($1,$2,$3) RETURNING id`
	var id int64
	if err := r.pool.QueryRow(context.Background(), q, p.Name, p.Price, p.Tags).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
