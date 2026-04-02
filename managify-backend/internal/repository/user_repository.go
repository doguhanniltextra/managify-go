package repository

import (
	"context"
	"managify/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserRepository, User veri işlemleri sözleşmesidir.
type UserRepository interface {
	FindOne(ctx context.Context, filter interface{}) (*models.User, error)
	InsertOne(ctx context.Context, document *models.User) error
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) error
	DeleteOne(ctx context.Context, filter interface{}) error

	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByFullName(ctx context.Context, fullName string) (*models.User, error)
	LinkGoogleID(ctx context.Context, userID interface{}, googleID string) error

	FindByID(ctx context.Context, id interface{}) (*models.User, error)
	FindByVerificationToken(ctx context.Context, token string) (*models.User, error)
	VerifyUser(ctx context.Context, userID interface{}) error

	FindAllUsers(ctx context.Context) ([]models.User, error)
	FindUsersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.User, error)
	DeleteByID(ctx context.Context, id interface{}) (int64, error)
	IncrementProjectSize(ctx context.Context, userID interface{}, amount int) error
}

type mongoUserRepository struct {
	*BaseRepository[models.User]
}


func NewUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{
		BaseRepository: NewBaseRepository[models.User](db, "users"),
	}
}

func (r *mongoUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	res, err := r.FindOne(ctx, bson.M{"email": email})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return res, err
}

func (r *mongoUserRepository) FindByFullName(ctx context.Context, fullName string) (*models.User, error) {
	res, err := r.FindOne(ctx, bson.M{"full_name": fullName})
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return res, err
}

func (r *mongoUserRepository) LinkGoogleID(ctx context.Context, userID interface{}, googleID string) error {
	return r.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"google_id": googleID}},
	)
}

func (r *mongoUserRepository) FindByID(ctx context.Context, id interface{}) (*models.User, error) {
	res, err := r.FindOne(ctx, bson.M{"_id": id})
	if err == mongo.ErrNoDocuments {
		return nil, nil // Error swallowing for missing ID commonly used in service logic
	}
	return res, err
}

func (r *mongoUserRepository) FindByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	return r.FindOne(ctx, bson.M{"verificationtoken": token})
}

func (r *mongoUserRepository) FindAllUsers(ctx context.Context) ([]models.User, error) {
	opts := options.Find().SetLimit(100).SetProjection(bson.M{"password": 0})
	cursor, err := r.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *mongoUserRepository) FindUsersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.User, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *mongoUserRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (r *mongoUserRepository) VerifyUser(ctx context.Context, userID interface{}) error {
	return r.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{"isverified": true, "verificationtoken": ""}},
	)
}

func (r *mongoUserRepository) IncrementProjectSize(ctx context.Context, userID interface{}, amount int) error {
	return r.UpdateOne(ctx,
		bson.M{"_id": userID},
		bson.M{"$inc": bson.M{"project_size": amount}},
	)
}
