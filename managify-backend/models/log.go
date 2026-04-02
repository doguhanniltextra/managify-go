package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"project_id" json:"project_id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Message   string             `bson:"message" json:"message"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}
