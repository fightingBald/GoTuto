package inmem

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports/outbound"
)

var (
	_ outbound.ProductRepository = (*InMemRepo)(nil)
	_ outbound.UserRepository    = (*InMemRepo)(nil)
	_ outbound.CommentRepository = (*InMemRepo)(nil)
)

// 简单的内存实现，用于本地开发/测试和示例 wiring
type InMemRepo struct {
	mu          sync.RWMutex
	products    map[int64]domain.Product
	nextProduct int64
	users       map[int64]domain.User
	comments    map[int64]domain.Comment
	nextComment int64
}

func NewInMemRepo() *InMemRepo {
	r := &InMemRepo{
		products:    make(map[int64]domain.Product),
		nextProduct: 1,
		users:       make(map[int64]domain.User),
		comments:    make(map[int64]domain.Comment),
		nextComment: 1,
	}
	// seed demo data
	r.products[1] = domain.Product{ID: 1, Name: "Blue Widget", Price: 1999}
	r.products[2] = domain.Product{ID: 2, Name: "Red Gizmo", Price: 2999}
	r.nextProduct = 3
	r.users[1] = domain.User{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: time.Date(2024, time.January, 10, 12, 0, 0, 0, time.UTC)}
	r.users[2] = domain.User{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: time.Date(2024, time.January, 11, 9, 30, 0, 0, time.UTC)}
	return r
}

func (r *InMemRepo) CreateComment(ctx context.Context, comment *domain.Comment) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextComment
	comment.ID = id
	if comment.CreatedAt.IsZero() {
		now := time.Now().UTC()
		comment.CreatedAt = now
		comment.UpdatedAt = now
	} else if comment.UpdatedAt.IsZero() {
		comment.UpdatedAt = comment.CreatedAt
	}
	r.comments[id] = *comment
	r.nextComment = id + 1
	return id, nil
}

func (r *InMemRepo) GetCommentByID(ctx context.Context, id int64) (*domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.comments[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copy := c
	return &copy, nil
}

func (r *InMemRepo) ListCommentsByProduct(ctx context.Context, productID int64) ([]domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []domain.Comment
	for _, c := range r.comments {
		if c.ProductID == productID {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].ID < out[j].ID
		}
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (r *InMemRepo) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.comments[comment.ID]; !ok {
		return domain.ErrNotFound
	}
	r.comments[comment.ID] = *comment
	return nil
}

func (r *InMemRepo) DeleteComment(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.comments[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.comments, id)
	return nil
}

func (r *InMemRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.products[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	// return copy
	pp := p
	return &pp, nil
}

func (r *InMemRepo) Search(ctx context.Context, q string, page, pageSize int) ([]domain.Product, int, error) {
	if page < 1 {
		page = 1
	}
	start := (page - 1) * pageSize
	q = strings.TrimSpace(strings.ToLower(q))

	r.mu.RLock()
	defer r.mu.RUnlock()
	var filtered []domain.Product
	for _, p := range r.products {
		if q == "" || strings.Contains(strings.ToLower(p.Name), q) {
			filtered = append(filtered, p)
		}
	}
	total := len(filtered)
	// simple pagination
	if start >= total {
		return []domain.Product{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

func (r *InMemRepo) Create(ctx context.Context, p *domain.Product) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.nextProduct
	p.ID = id
	r.products[id] = *p
	r.nextProduct = id + 1
	return id, nil
}

func (r *InMemRepo) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.products, id)
	return nil
}

func (r *InMemRepo) Update(ctx context.Context, p *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.products[p.ID]; !ok {
		return domain.ErrNotFound
	}
	r.products[p.ID] = *p
	return nil
}

func (r *InMemRepo) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	uu := u
	return &uu, nil
}
