package httpadapter

import (
	"net/http"

	httpadp "github.com/fightingBald/GoTuto/internal/adapters/http"
)

// 已在 handler_product_read.go 中实现 SearchProducts，此文件提供同样的签名以匹配期望的命名
//（保留以便将来按文件分离）

func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request, params httpadp.SearchProductsParams) {
	// 直接复用 SearchProducts 的实现
	s.SearchProducts(w, r, params)
}
