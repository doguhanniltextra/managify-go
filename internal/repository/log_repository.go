package repository

import (
	"context"
	"managify/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LogRepository interface {
	InsertOne(ctx context.Context, log *models.ProjectLog) error
	FindByProjectID(ctx context.Context, projectID string) ([]models.ProjectLog, error)
	GetRecentUserLogs(ctx context.Context, userID string, limit int64) ([]models.ProjectLog, error)
}

type mongoLogRepository struct {
	*BaseRepository[models.ProjectLog]
}

func NewLogRepository(db *mongo.Database) LogRepository {
	return &mongoLogRepository{
		BaseRepository: NewBaseRepository[models.ProjectLog](db, "logs"),
	}
}

func (r *mongoLogRepository) InsertOne(ctx context.Context, log *models.ProjectLog) error {
	_, err := r.Collection.InsertOne(ctx, log)
	return err
}

func (r *mongoLogRepository) FindByProjectID(ctx context.Context, projectID string) ([]models.ProjectLog, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.ProjectLog
	for cursor.Next(ctx) {
		var logEntry models.ProjectLog
		if err := cursor.Decode(&logEntry); err != nil {
			continue
		}
		logs = append(logs, logEntry)
	}

	if logs == nil {
		logs = []models.ProjectLog{}
	}
	return logs, nil
}

func (r *mongoLogRepository) GetRecentUserLogs(ctx context.Context, userID string, limit int64) ([]models.ProjectLog, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(limit)
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.ProjectLog
	for cursor.Next(ctx) {
		var logEntry models.ProjectLog
		if err := cursor.Decode(&logEntry); err != nil {
			continue
		}
		logs = append(logs, logEntry)
	}

	if logs == nil {
		logs = []models.ProjectLog{}
	}
	return logs, nil
}
