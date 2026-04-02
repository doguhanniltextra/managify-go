package repository

import (
	"context"
	"managify/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectInviteRepository interface {
	CountByFilter(ctx context.Context, receiverID interface{}, projectID interface{}, statuses []string) (int64, error)
	UpsertInvite(ctx context.Context, receiverID, projectID, senderID interface{}) (*models.ProjectInvite, error)
	FindInvitesByReceiverID(ctx context.Context, receiverID interface{}) ([]*models.ProjectInvite, error)
	UpdateStatus(ctx context.Context, inviteID, receiverID interface{}, newStatus string) (*models.ProjectInvite, error)
}

type mongoProjectInviteRepository struct {
	*BaseRepository[models.ProjectInvite]
}

func NewProjectInviteRepository(db *mongo.Database) ProjectInviteRepository {
	return &mongoProjectInviteRepository{
		BaseRepository: NewBaseRepository[models.ProjectInvite](db, "project_invites"),
	}
}

func (r *mongoProjectInviteRepository) CountByFilter(ctx context.Context, receiverID interface{}, projectID interface{}, statuses []string) (int64, error) {
	filter := bson.M{
		"receiver_id": receiverID,
		"project_id":  projectID,
		"status":      bson.M{"$in": statuses},
	}
	return r.Collection.CountDocuments(ctx, filter)
}

func (r *mongoProjectInviteRepository) UpsertInvite(ctx context.Context, receiverID, projectID, senderID interface{}) (*models.ProjectInvite, error) {
	statusFilter := bson.M{"$in": []string{"pending", "accepted"}}
	filter := bson.M{
		"receiver_id": receiverID,
		"project_id":  projectID,
		"status":      statusFilter,
	}

	update := bson.M{
		"$setOnInsert": bson.M{
			"project_id":  projectID,
			"receiver_id": receiverID,
			"sender_id":   senderID,
			"status":      "pending",
			"created_at":  time.Now(),
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	res := r.Collection.FindOneAndUpdate(ctx, filter, update, opts)

	var invite models.ProjectInvite
	if err := res.Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}

func (r *mongoProjectInviteRepository) FindInvitesByReceiverID(ctx context.Context, receiverID interface{}) ([]*models.ProjectInvite, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"receiver_id": receiverID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invites []*models.ProjectInvite
	for cursor.Next(ctx) {
		var invite models.ProjectInvite
		if err := cursor.Decode(&invite); err != nil {
			return nil, err
		}
		invites = append(invites, &invite)
	}

	if invites == nil {
		invites = []*models.ProjectInvite{}
	}
	return invites, nil
}

func (r *mongoProjectInviteRepository) UpdateStatus(ctx context.Context, inviteID, receiverID interface{}, newStatus string) (*models.ProjectInvite, error) {
	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}
	res := r.Collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": inviteID, "receiver_id": receiverID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var invite models.ProjectInvite
	if err := res.Decode(&invite); err != nil {
		return nil, err
	}
	return &invite, nil
}
