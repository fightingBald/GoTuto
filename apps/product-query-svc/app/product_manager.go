package app

import (
	"context"

	appsports "github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

// ProductManager 可用于事务/上下文/跨服务协调（占位实现）
type ProductManager struct{ svc appsports.ProductQueryPort }

func NewProductManager(svc appsports.ProductQueryPort) *ProductManager {
	return &ProductManager{svc: svc}
}

func (m *ProductManager) GetProduct(ctx context.Context, id int64) (interface{}, error) {
	// 目前只是简单代理到应用服务
	return m.svc.GetProduct(id)
}
