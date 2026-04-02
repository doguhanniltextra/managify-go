package repository

import (
	"context"
	"managify/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectRepository interface {
	FindAll(ctx context.Context) ([]models.Project, error)
	InsertOne(ctx context.Context, project *models.Project) error
	DeleteByID(ctx context.Context, id interface{}) (int64, error)
	FindOneWithAccess(ctx context.Context, projectID interface{}, userID interface{}) (*models.Project, error)
	FindAllByUserID(ctx context.Context, userID interface{}) ([]*models.Project, error)
	CheckUserInProject(ctx context.Context, projectID interface{}, userID interface{}) (bool, error)
	RemoveMemberFromProject(ctx context.Context, ownerID interface{}, memberID interface{}) (int64, error)
	VerifyProject(ctx context.Context, projectID interface{}) (bool, error)
	VerifyProjectOwner(ctx context.Context, projectID interface{}, ownerID interface{}) (bool, error)
	AddUserToProject(ctx context.Context, projectID interface{}, userID interface{}) error
}

type mongoProjectRepository struct {
	*BaseRepository[models.Project]
}

func NewProjectRepository(db *mongo.Database) ProjectRepository {
	return &mongoProjectRepository{
		BaseRepository: NewBaseRepository[models.Project](db, "projects"),
	}
}

func (r *mongoProjectRepository) FindAll(ctx context.Context) ([]models.Project, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *mongoProjectRepository) InsertOne(ctx context.Context, project *models.Project) error {
	_, err := r.Collection.InsertOne(ctx, project)
	return err
}

func (r *mongoProjectRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *mongoProjectRepository) FindOneWithAccess(ctx context.Context, projectID interface{}, userID interface{}) (*models.Project, error) {
	var project models.Project
	err := r.Collection.FindOne(ctx, bson.M{
		"_id": projectID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"team": bson.M{"$in": []interface{}{userID}}},
		},
	}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&project)

	if err == mongo.ErrNoDocuments {
		return nil, nil // return nil gracefully, handling is above in service
	}
	return &project, err
}

func (r *mongoProjectRepository) FindAllByUserID(ctx context.Context, userID interface{}) ([]*models.Project, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"owner_id": userID},
			{"team": bson.M{"$in": []interface{}{userID}}},
		},
	}
	opts := options.Find().SetLimit(50)
	cursor, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*models.Project
	if err = cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	if projects == nil { // return empty array instead of nil
		projects = []*models.Project{}
	}
	return projects, nil
}

func (r *mongoProjectRepository) CheckUserInProject(ctx context.Context, projectID interface{}, userID interface{}) (bool, error) {
	filter := bson.M{
		"_id": projectID,
		"$or": []bson.M{
			{"owner_id": userID},
			{"team": userID},
		},
	}
	count, err := r.Collection.CountDocuments(ctx, filter)
	return count > 0, err
}

func (r *mongoProjectRepository) RemoveMemberFromProject(ctx context.Context, ownerID interface{}, memberID interface{}) (int64, error) {
	res, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"owner_id": ownerID},
		bson.M{"$pull": bson.M{"team": memberID}},
	)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (r *mongoProjectRepository) VerifyProject(ctx context.Context, projectID interface{}) (bool, error) {
	var project models.Project
	err := r.Collection.FindOne(ctx, bson.M{"_id": projectID}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *mongoProjectRepository) VerifyProjectOwner(ctx context.Context, projectID interface{}, ownerID interface{}) (bool, error) {
	var project models.Project
	err := r.Collection.FindOne(ctx, bson.M{
		"_id":      projectID,
		"owner_id": ownerID,
	}).Decode(&project)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *mongoProjectRepository) AddUserToProject(ctx context.Context, projectID interface{}, userID interface{}) error {
	update := bson.M{
		"$addToSet": bson.M{"team": userID},
	}
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": projectID}, update)
	return err
}
