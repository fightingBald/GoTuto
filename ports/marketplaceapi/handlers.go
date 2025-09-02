package marketplaceapi

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
)

type Server struct {
	products map[int64]Product
}

func NewServer() *Server {
	return &Server{
		products: map[int64]Product{
			1: {Id: 1, Name: "Apple", Price: 1.23},
			2: {Id: 2, Name: "Pen", Price: 2.50},
			3: {Id: 3, Name: "Pencil", Price: 0.99},
		},
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

// operationId: GetProductByID
func (s *Server) GetProductByID(w http.ResponseWriter, r *http.Request, id int64) {
	p, ok := s.products[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, Error{Code: "NOT_FOUND", Message: "product not found"})
		return
	}
	writeJSON(w, http.StatusOK, p) // 直接回 Product（你的 spec 就是这样定义的）
}

// operationId: SearchProducts
func (s *Server) SearchProducts(w http.ResponseWriter, r *http.Request, params SearchProductsParams) {
	q := ""
	if params.Q != nil {
		q = *params.Q
	}
	page := 1
	if params.Page != nil && *params.Page > 0 {
		page = *params.Page
	}
	pageSize := 10
	if params.PageSize != nil && *params.PageSize > 0 {
		pageSize = *params.PageSize
	}

	// 过滤 + 稳定排序
	var all []Product
	for _, p := range s.products {
		if q == "" || strings.Contains(strings.ToLower(p.Name), strings.ToLower(q)) {
			all = append(all, p)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].Id < all[j].Id })

	total := len(all)
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	writeJSON(w, http.StatusOK, ProductList{
		Items:    all[start:end],
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}
