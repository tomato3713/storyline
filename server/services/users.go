package services

import (
	"context"
	"log/slog"

	"github.com/tomato3713/storyline/server/model"
	"github.com/tomato3713/storyline/server/repository"
)

type userService struct {
	exec repository.Repository
}

func (u *userService) CreateUserByName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{
		Name: name,
	}
	createdUser, err := u.exec.InsertUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *userService) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	user, err := u.exec.GetUserByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userService) GetUserAll(ctx context.Context) ([]*model.User, error) {
	slog.Info("begin userService.GetUserAll()")
	users, err := u.exec.GetUserAll(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("finished userService.GetUserAll")
	return users, nil
}
