package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrEmailValidation = errors.New("email validation error")
)

func IsValidEmail(email string) bool {
	// 这是一个常用的简化正则，够用 90% 场景
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name string, email string) (*User, error) {
	u := &User{Name: name, Email: email, CreatedAt: time.Now()}

	if err := u.Validate(); err != nil {
		return nil, ErrValidation
	}

	return u, nil
}

func (u *User) Validate() error {
	if IsValidEmail(u.Email) {
		return ErrValidation
	}
	return nil
}

//TODO how to avoid a same email could create several account

func (u *User) ChangeName(newName string) error {
	u.Name = newName
	return nil
}
