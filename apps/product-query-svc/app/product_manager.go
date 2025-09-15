package app

import (
	"github.com/fightingBald/GoTuto/internal/domain"
	"github.com/fightingBald/GoTuto/internal/ports"
)

// ProductManager 可用于事务/上下文/跨服务协调（占位实现）
type ProductManager struct{ svc ports.ProductService }

func NewProductManager(svc ports.ProductService) *ProductManager {
	return &ProductManager{svc: svc}
}

func (m *ProductManager) GetProduct(id int64) (*domain.Product, error) {
	// 目前只是简单代理到应用服务
	return m.svc.GetProduct(id)
}
