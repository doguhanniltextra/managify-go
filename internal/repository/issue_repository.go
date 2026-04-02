package repository

import (
	"context"
	"managify/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IssueRepository interface {
	InsertOne(ctx context.Context, issue *models.Issue) error
	DeleteByID(ctx context.Context, id interface{}) (int64, error)
	FindByID(ctx context.Context, id interface{}) (*models.Issue, error)
	FindByStatusID(ctx context.Context, statusID interface{}) ([]*models.Issue, error)
	UpdateStatus(ctx context.Context, issueID interface{}, newStatusID interface{}) (int64, error)
	FindOncomingIssues(ctx context.Context, projectID interface{}, currentTime time.Time, limitTime time.Time) ([]*models.Issue, error)
}

type mongoIssueRepository struct {
	*BaseRepository[models.Issue]
}

func NewIssueRepository(db *mongo.Database) IssueRepository {
	return &mongoIssueRepository{
		BaseRepository: NewBaseRepository[models.Issue](db, "issues"),
	}
}

func (r *mongoIssueRepository) InsertOne(ctx context.Context, issue *models.Issue) error {
	_, err := r.Collection.InsertOne(ctx, issue)
	return err
}

func (r *mongoIssueRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *mongoIssueRepository) FindByID(ctx context.Context, id interface{}) (*models.Issue, error) {
	var issue models.Issue
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&issue)
	if err == mongo.ErrNoDocuments {
		return nil, nil // return nil for standard repo behavior
	}
	return &issue, err
}

func (r *mongoIssueRepository) FindByStatusID(ctx context.Context, statusID interface{}) ([]*models.Issue, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"status_id": statusID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var issues []*models.Issue
	if err = cursor.All(ctx, &issues); err != nil {
		return nil, err
	}
	if issues == nil {
		issues = []*models.Issue{}
	}
	return issues, nil
}

func (r *mongoIssueRepository) UpdateStatus(ctx context.Context, issueID interface{}, newStatusID interface{}) (int64, error) {
	update := bson.M{
		"$set": bson.M{
			"status_id":  newStatusID,
			"updated_at": time.Now(),
		},
	}
	res, err := r.Collection.UpdateOne(ctx, bson.M{"_id": issueID}, update)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *mongoIssueRepository) FindOncomingIssues(ctx context.Context, projectID interface{}, currentTime time.Time, limitTime time.Time) ([]*models.Issue, error) {
	format := "2006-01-02"

	filter := bson.M{
		"project_id": projectID,
		"due_date": bson.M{
			"$gte": currentTime.Format(format),
			"$lte": limitTime.Format(format),
		},
	}

	projection := bson.M{
		"title":       1,
		"description": 1,
		"due_date":    1,
	}
	opt := options.Find().SetSort(bson.D{{Key: "due_date", Value: 1}}).SetProjection(projection)

	cursor, err := r.Collection.Find(ctx, filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var issues []*models.Issue
	if err = cursor.All(ctx, &issues); err != nil {
		return nil, err
	}
	if issues == nil {
		issues = []*models.Issue{}
	}
	return issues, nil
}
