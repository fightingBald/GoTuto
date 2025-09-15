package postgres

import (
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
	internaldb "github.com/fightingBald/GoTuto/internal/adapters/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewProductRepository(pool *pgxpool.Pool) ports.ProductRepository {
	return internaldb.NewPGProductRepo(pool)
}
