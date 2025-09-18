package domain

import "errors"

// 领域错误（供适配器映射状态码）
var (
    ErrValidation = errors.New("validation error")
    ErrNotFound   = errors.New("not found")
)

type Product struct {
    ID    int64
    Name  string
    Price int64   // 分为单位（避免浮点）
    Tags  []string
}

// 工厂方法：统一创建入口，保证不变式
func NewProduct(name string, priceCents int64, tags []string) (*Product, error) {
    p := &Product{Name: name, Price: priceCents}
    if len(tags) > 0 {
        // 去重
        seen := map[string]struct{}{}
        for _, t := range tags {
            if t == "" { continue }
            if _, ok := seen[t]; ok { continue }
            seen[t] = struct{}{}
            p.Tags = append(p.Tags, t)
        }
    }
    if err := p.Validate(); err != nil {
        return nil, err
    }
    return p, nil
}

// 基础不变式校验
func (p *Product) Validate() error {
    if p.Name == "" {
        return ErrValidation
    }
    if p.Price < 0 {
        return ErrValidation
    }
    return nil
}

// 富领域行为示例：修改价格（非负）
func (p *Product) ChangePrice(newPriceCents int64) error {
    if newPriceCents < 0 {
        return ErrValidation
    }
    p.Price = newPriceCents
    return nil
}

// 富领域行为示例：添加标签（去重，最多 5 个）
func (p *Product) AddTag(tag string) error {
    if tag == "" { return nil }
    for _, t := range p.Tags {
        if t == tag { return nil }
    }
    if len(p.Tags) >= 5 {
        return ErrValidation
    }
    p.Tags = append(p.Tags, tag)
    return nil
}

// 富领域行为示例：移除标签
func (p *Product) RemoveTag(tag string) {
    out := p.Tags[:0]
    for _, t := range p.Tags {
        if t != tag { out = append(out, t) }
    }
    p.Tags = out
}
