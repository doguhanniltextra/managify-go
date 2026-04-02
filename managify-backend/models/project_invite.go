package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectInvite struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID  primitive.ObjectID `bson:"project_id"`
	SenderID   primitive.ObjectID `bson:"sender_id"`
	ReceiverID primitive.ObjectID `bson:"receiver_id"`
	Status     string             `bson:"status"` // "pending", "accepted", "declined"
	CreatedAt  time.Time          `bson:"created_at"`
}
