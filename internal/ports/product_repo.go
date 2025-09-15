package ports

import "github.com/fightingBald/GoTuto/internal/domain"

type ProductRepo interface {
	GetByID(id int64) (*domain.Product, error)
	Search(q string, page, pageSize int) ([]domain.Product, error)
	Create(p *domain.Product) (int64, error)
}
