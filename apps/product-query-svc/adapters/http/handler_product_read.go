package httpadapter

import (
	"encoding/json"
	"net/http"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

type Server struct{ svc ports.ProductService }

func NewServer(s ports.ProductService) *Server { return &Server{svc: s} }

func (s *Server) GetProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	p, err := s.svc.GetProduct(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(Error{Code: "NOT_FOUND", Message: err.Error()})
		return
	}
	out := Product{Id: p.ID, Name: p.Name, Price: float32(p.Price) / 100.0}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

func (s *Server) SearchProducts(w http.ResponseWriter, r *http.Request, params SearchProductsParams) {
	q := ""
	if params.Q != nil {
		q = *params.Q
	}
	page := 1
	if params.Page != nil {
		page = *params.Page
	}
	pageSize := 20
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}
	items, err := s.svc.SearchProducts(q, page, pageSize)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(Error{Code: "INTERNAL", Message: err.Error()})
		return
	}
	var out []Product
	for _, it := range items {
		out = append(out, Product{Id: it.ID, Name: it.Name, Price: float32(it.Price) / 100.0})
	}
	resp := ProductList{Items: out, Page: page, PageSize: pageSize, Total: len(out)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Healthz 健康检查
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
