package repository

import (
	"context"
	"managify/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleRepository interface {
	FindAll(ctx context.Context) ([]models.Role, error)
	InsertOne(ctx context.Context, role *models.Role) error
	DeleteByID(ctx context.Context, id interface{}) (int64, error)
}

type mongoRoleRepository struct {
	*BaseRepository[models.Role]
}

func NewRoleRepository(db *mongo.Database) RoleRepository {
	return &mongoRoleRepository{
		BaseRepository: NewBaseRepository[models.Role](db, "roles"),
	}
}

func (r *mongoRoleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []models.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *mongoRoleRepository) InsertOne(ctx context.Context, role *models.Role) error {
	_, err := r.Collection.InsertOne(ctx, role)
	return err
}

func (r *mongoRoleRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
