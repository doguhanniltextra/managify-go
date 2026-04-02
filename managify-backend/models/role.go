package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	ProjectID primitive.ObjectID `bson:"project_id" json:"project_id"`
	RoleName  string             `bson:"role" json:"role"`
}
