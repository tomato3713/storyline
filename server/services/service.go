package services

import (
	"context"

	"github.com/tomato3713/storyline/server/model"
	"github.com/tomato3713/storyline/server/repository"
)

type Services interface {
	UserService
}

type UserService interface {
	CreateUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserAll(ctx context.Context) ([]*model.User, error)
}

type services struct {
	*userService
}

func New(exec repository.Repository) Services {
	return &services{
		userService: &userService{exec: exec},
	}
}
