package domain

import (
	"regexp"
	"strings"
	"time"
)

var emailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	return emailRegexp.MatchString(strings.TrimSpace(email))
}

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name string, email string) (*User, error) {
	u := &User{
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(email),
		CreatedAt: time.Now().UTC(),
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errValidation("name required")
	}
	if !IsValidEmail(u.Email) {
		return errValidation("invalid email format")
	}
	return nil
}

//TODO how to avoid a same email could create several account

func (u *User) ChangeName(newName string) error {
	cleaned := strings.TrimSpace(newName)
	if cleaned == "" {
		return errValidation("name required")
	}
	u.Name = cleaned
	return nil
}
