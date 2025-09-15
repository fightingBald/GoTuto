package app

import (
	appsports "github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
	"github.com/fightingBald/GoTuto/internal"
)

// NewProductQueryService 返回一个实现入站端口的应用服务，内部复用 existing internal.NewProductService
func NewProductQueryService(repo appsports.ProductRepository) appsports.ProductQueryPort {
	return internal.NewProductService(repo)
}
