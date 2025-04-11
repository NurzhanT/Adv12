package repo

import (
	"context"
	"errors"
	"inventory/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type productMongoRepo struct {
	collection *mongo.Collection
}

func NewProductMongoRepo(client *mongo.Client, dbName string) domain.ProductRepository {
	collection := client.Database(dbName).Collection("products")
	return &productMongoRepo{collection: collection}
}

func (r *productMongoRepo) Create(p domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, p)
	return err
}

func (r *productMongoRepo) GetByID(id string) (domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var p domain.Product
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&p)
	if err != nil {
		return domain.Product{}, errors.New("not found")
	}
	return p, nil
}

func (r *productMongoRepo) Update(id string, p domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": p})
	return err
}

func (r *productMongoRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *productMongoRepo) List() ([]domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []domain.Product
	for cursor.Next(ctx) {
		var p domain.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
