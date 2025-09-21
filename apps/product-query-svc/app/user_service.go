package app

import (
	"context"
	"errors"

	"github.com/fightingBald/GoTuto/apps/product-query-svc/domain"
	"github.com/fightingBald/GoTuto/apps/product-query-svc/ports"
)

type UserService struct {
	repo ports.UserRepo
}

func NewUserService(r ports.UserRepo) *UserService { return &UserService{repo: r} }

func (s *UserService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.Join(domain.ErrValidation, errors.New("id must be a positive integer"))
	}
	return s.repo.GetUserByID(ctx, id)
}
