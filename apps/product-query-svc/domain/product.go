package domain

import (
	"errors"
	"strings"
)

// 领域错误（供适配器映射状态码）
var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

// Product 是领域聚合根，Price 以分为单位避免浮点误差。
type Product struct {
	ID    int64
	Name  string
	Price int64
	Tags  []string
}

const maxTags = 5

// NewProduct 统一入口，构建并校验不变式。
func NewProduct(name string, priceCents int64, tags []string) (*Product, error) {
	p := &Product{Name: strings.TrimSpace(name), Price: priceCents}
	if err := p.replaceTags(tags); err != nil {
		return nil, err
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

// Validate 检查核心不变式。
func (p *Product) Validate() error {
	if p.Name == "" {
		return errValidation("name required")
	}
	if p.Price < 0 {
		return errValidation("price must be >= 0")
	}
	if len(p.Tags) > maxTags {
		return errValidation("tags exceed limit")
	}
	return nil
}

// ChangePrice 变更价格（分为单位）。
func (p *Product) ChangePrice(newPrice int64) error {
	if newPrice < 0 {
		return errValidation("price must be >= 0")
	}
	p.Price = newPrice
	return nil
}

// AddTag 添加标签，自动去重并限制数量。
func (p *Product) AddTag(tag string) error {
	cleaned := strings.TrimSpace(tag)
	if cleaned == "" {
		return nil
	}
	for _, t := range p.Tags {
		if equalFold(t, cleaned) {
			return nil
		}
	}
	if len(p.Tags) >= maxTags {
		return errValidation("tags exceed limit")
	}
	p.Tags = append(p.Tags, cleaned)
	return nil
}

// RemoveTag 移除指定标签（按不区分大小写匹配）。
func (p *Product) RemoveTag(tag string) {
	cleaned := strings.TrimSpace(tag)
	if cleaned == "" {
		return
	}
	out := p.Tags[:0]
	for _, t := range p.Tags {
		if !equalFold(t, cleaned) {
			out = append(out, t)
		}
	}
	p.Tags = out
}

// replaceTags 重建标签列表（调用方负责去重复构建）。
func (p *Product) replaceTags(tags []string) error {
	sanitized, err := sanitizeTags(tags)
	if err != nil {
		return err
	}
	p.Tags = sanitized
	return nil
}

// errValidation 构造带细节的校验错误。
func errValidation(msg string) error {
	return errors.Join(ErrValidation, errors.New(msg))
}

func equalFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

func sanitizeTags(tags []string) ([]string, error) {
	if len(tags) == 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, len(tags))
	sanitized := make([]string, 0, len(tags))
	for _, raw := range tags {
		cleaned := strings.TrimSpace(raw)
		if cleaned == "" {
			continue
		}
		key := strings.ToLower(cleaned)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		sanitized = append(sanitized, cleaned)
		if len(sanitized) > maxTags {
			return nil, errValidation("tags exceed limit")
		}
	}
	if len(sanitized) == 0 {
		return nil, nil
	}
	return sanitized, nil
}
