//go:build ignore
// +build ignore

package httpadp

// NOTE: 该文件被标记为被构建忽略，保留为历史/参考。实际运行使用 apps/ 下的 adapters/http 实现。

/*
package httpadp

import (
	"encoding/json"
	"net/http"

	"github.com/fightingBald/GoTuto/internal/domain"
	"github.com/fightingBald/GoTuto/internal/ports"
)

type Server struct{ svc ports.ProductService }

func NewServer(s ports.ProductService) *Server { return &Server{svc: s} }

// helper: convert domain.Product -> generated Product
func toProduct(p domain.Product) Product {
	return Product{
		Id:    p.ID,
		Name:  p.Name,
		Price: float32(p.Price) / 100.0,
	}
}

// 假设你的 openapi 有 /products/{id} 和 /products/search
func (s *Server) GetProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	p, err := s.svc.GetProduct(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(Error{Code: "NOT_FOUND", Message: err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toProduct(*p))
}

func (s *Server) SearchProducts(w http.ResponseWriter, r *http.Request, params SearchProductsParams) {
	page := 1
	if params.Page != nil {
		page = *params.Page
	}
	pageSize := 20
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	q := ""
	if params.Q != nil {
		q = *params.Q
	}

	items, err := s.svc.SearchProducts(q, page, pageSize)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Code: "INTERNAL", Message: err.Error()})
		return
	}
	// convert
	var out []Product
	for _, it := range items {
		out = append(out, toProduct(it))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(ProductList{Items: out, Page: page, PageSize: pageSize, Total: len(out)})
}

// 可选：健康检查（非 openapi 路由）
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("ok"))
}
*/
