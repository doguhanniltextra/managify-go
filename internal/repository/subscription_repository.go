package repository

import (
	"context"
	"managify/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository interface {
	FindOne(ctx context.Context, filter interface{}) (*models.Subscription, error)
	InsertOne(ctx context.Context, document *models.Subscription) error

	FindByUserID(ctx context.Context, userID interface{}) (*models.Subscription, error)
}

type mongoSubscriptionRepository struct {
	*BaseRepository[models.Subscription]
}

func NewSubscriptionRepository(db *mongo.Database) SubscriptionRepository {
	return &mongoSubscriptionRepository{
		BaseRepository: NewBaseRepository[models.Subscription](db, "subscriptions"),
	}
}

func (r *mongoSubscriptionRepository) FindByUserID(ctx context.Context, userID interface{}) (*models.Subscription, error) {
	res, err := r.FindOne(ctx, bson.M{"user_id": userID})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return res, err
}
