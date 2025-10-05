package httpadapter

import (
	"math"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func presentProduct(p *domain.Product) Product {
	if p == nil {
		return Product{}
	}
	return Product{
		Id:    p.ID,
		Name:  p.Name,
		Price: centsToAmount(p.Price),
	}
}

func presentProducts(items []domain.Product) []Product {
	if len(items) == 0 {
		return []Product{}
	}
	out := make([]Product, 0, len(items))
	for i := range items {
		out = append(out, presentProduct(&items[i]))
	}
	return out
}

func centsToAmount(cents int64) float32 {
	return float32(cents) / 100.0
}

func amountToCents(amount float32) int64 {
	return int64(math.Round(float64(amount) * 100.0))
}

func presentUser(u *domain.User) User {
	if u == nil {
		return User{}
	}
	id := u.ID
	createdAt := u.CreatedAt.UTC()
	return User{
		Id:        &id,
		Name:      u.Name,
		Email:     openapi_types.Email(u.Email),
		CreatedAt: &createdAt,
	}
}

func presentComment(c *domain.Comment) Comment {
	if c == nil {
		return Comment{}
	}
	return Comment{
		Id:        c.ID,
		ProductId: c.ProductID,
		UserId:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt.UTC(),
		UpdatedAt: c.UpdatedAt.UTC(),
	}
}

func presentComments(items []domain.Comment) []Comment {
	if len(items) == 0 {
		return []Comment{}
	}
	out := make([]Comment, 0, len(items))
	for i := range items {
		out = append(out, presentComment(&items[i]))
	}
	return out
}
