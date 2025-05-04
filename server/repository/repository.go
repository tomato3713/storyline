package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/tomato3713/storyline/server/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	InsertUser(context.Context, *model.User) (*model.User, error)
	GetUserByName(ctx context.Context, name string) (*model.User, error)
	GetUserAll(ctx context.Context) ([]*model.User, error)
}

type repositoryImpl struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *repositoryImpl {
	return &repositoryImpl{
		db: db,
	}
}

type user struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedAt time.Time
}

// GetUsers は、データベースからユーザーを取得するメソッドです。
func (r *repositoryImpl) GetUserByName(ctx context.Context, name string) (*model.User, error) {
	var user user
	err := r.db.Collection("users").FindOne(ctx, bson.D{{Key: "name", Value: name}}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:   user.ID.Hex(),
		Name: user.Name,
	}, nil
}

func (r *repositoryImpl) GetUserAll(ctx context.Context) ([]*model.User, error) {
	findOptions := options.Find()
	findOptions.SetLimit(100)
	cur, err := r.db.Collection("users").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	var users []user
	for cur.Next(ctx) {
		var elem user
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}

		users = append(users, elem)
	}

	results := make([]*model.User, len(users))
	for i, user := range users {
		mu := model.User{
			ID:   user.ID.Hex(),
			Name: user.Name,
		}
		results[i] = &mu
	}

	return results, nil
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
