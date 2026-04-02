package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type BaseRepository[T any] struct {
	Collection *mongo.Collection
}

func NewBaseRepository[T any](db *mongo.Database, collectionName string) *BaseRepository[T] {
	return &BaseRepository[T]{
		Collection: db.Collection(collectionName),
	}
}
func (r *BaseRepository[T]) InsertOne(ctx context.Context, document *T) error {
	_, err := r.Collection.InsertOne(ctx, document)
	return err
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, filter interface{}) (*T, error) {
	var result T
	err := r.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *BaseRepository[T]) UpdateOne(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *BaseRepository[T]) DeleteOne(ctx context.Context, filter interface{}) error {
	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}
