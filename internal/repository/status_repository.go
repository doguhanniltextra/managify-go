package repository

import (
	"context"
	"managify/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatusRepository interface {
	InsertOne(ctx context.Context, status *models.Status) error
	DeleteByID(ctx context.Context, id interface{}) (int64, error)
	FindByProjectID(ctx context.Context, projectID interface{}) ([]*models.Status, error)
}

type mongoStatusRepository struct {
	*BaseRepository[models.Status]
}

func NewStatusRepository(db *mongo.Database) StatusRepository {
	return &mongoStatusRepository{
		BaseRepository: NewBaseRepository[models.Status](db, "status"),
	}
}

func (r *mongoStatusRepository) InsertOne(ctx context.Context, status *models.Status) error {
	_, err := r.Collection.InsertOne(ctx, status)
	return err
}

func (r *mongoStatusRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *mongoStatusRepository) FindByProjectID(ctx context.Context, projectID interface{}) ([]*models.Status, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var statuses []*models.Status
	if err := cursor.All(ctx, &statuses); err != nil {
		return nil, err
	}
	if statuses == nil {
		statuses = []*models.Status{}
	}
	return statuses, nil
}
