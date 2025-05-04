package repository

import (
	"context"
	"fmt"

	"github.com/tomato3713/storyline/server/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertUser(context.Context, model.User) error
}

type repositoryImpl struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repositoryImpl {
	return &repositoryImpl{
		db: db,
	}
}

// InsertUser は、ユーザーをデータベースに挿入するメソッドです。
func (r *repositoryImpl) InsertUser(ctx context.Context, user model.User) error {
	_, err := r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}
