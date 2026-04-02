package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	ProjectID primitive.ObjectID   `bson:"project_id" json:"project_id"`
	CreatorID primitive.ObjectID   `bson:"creator_id" json:"-"`
	Name      string               `bson:"name" json:"name"`
	IssueIDs  []primitive.ObjectID `bson:"issues" json:"issues_id"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
