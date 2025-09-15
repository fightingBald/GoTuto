package postgres

import (
	internaldb "github.com/fightingBald/GoTuto/internal/adapters/db"
	"github.com/fightingBald/GoTuto/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewProductRepository(pool *pgxpool.Pool) ports.ProductRepo {
	return internaldb.NewPGProductRepo(pool)
}
