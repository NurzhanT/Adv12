package repo

import (
	"context"
	"errors"
	"inventory/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepo struct {
	collection *mongo.Collection
}

func NewUserMongoRepo(client *mongo.Client, dbName string) domain.UserRepository {
	return &userMongoRepo{
		collection: client.Database(dbName).Collection("users"),
	}
}

func (r *userMongoRepo) Create(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check existing user
	existing, _ := r.GetByUsername(user.Username)
	if existing.Username != "" {
		return errors.New("username already exists")
	}

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userMongoRepo) GetByUsername(username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}
	return user, nil
}
