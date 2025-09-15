package domain

type Product struct {
	ID    int64
	Name  string
	Price int64 // 分为单位，别玩浮点
	Tags  []string
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return Err("name required")
	}
	if p.Price < 0 {
		return Err("price must be >= 0")
	}
	return nil
}

type domainErr string

func (e domainErr) Error() string { return string(e) }
func Err(msg string) error        { return domainErr(msg) }
