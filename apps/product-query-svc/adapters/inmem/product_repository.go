package inmem

import (
	"errors"
	"strings"
	"sync"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

// 简单的内存实现，用于本地开发/测试和示例 wiring
type InMemRepo struct {
	mu   sync.RWMutex
	data map[int64]domain.Product
	next int64
}

func NewInMemRepo() ports.ProductRepo {
	r := &InMemRepo{data: make(map[int64]domain.Product), next: 1}
	// seed demo data
	r.data[1] = domain.Product{ID: 1, Name: "Blue Widget", Price: 1999}
	r.data[2] = domain.Product{ID: 2, Name: "Red Gizmo", Price: 2999}
	r.next = 3
	return r
}

func (r *InMemRepo) GetByID(id int64) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	// return copy
	pp := p
	return &pp, nil
}

func (r *InMemRepo) Search(q string, page, pageSize int) ([]domain.Product, error) {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	q = strings.TrimSpace(strings.ToLower(q))

	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Product
	for _, p := range r.data {
		if q == "" || strings.Contains(strings.ToLower(p.Name), q) {
			out = append(out, p)
		}
	}
	// simple pagination
	if start >= len(out) {
		return []domain.Product{}, nil
	}
	end := start + pageSize
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], nil
}

func (r *InMemRepo) Create(p *domain.Product) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.next
	p.ID = id
	r.data[id] = *p
	r.next = id + 1
	return id, nil
}
