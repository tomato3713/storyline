package repository

import (
	"context"
	"fmt"

	"github.com/tomato3713/storyline/server/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	InsertUser(context.Context, *model.User) (*model.User, error)
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
func (r *repositoryImpl) InsertUser(ctx context.Context, user *model.User) (*model.User, error) {
	result, err := r.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}

	return &model.User{
		ID:   oid.Hex(),
		Name: user.Name,
	}, nil
}
