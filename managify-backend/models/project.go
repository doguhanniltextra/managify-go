package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	Category    string               `bson:"category" json:"category"`
	Tags        []string             `bson:"tags,omitempty" json:"tags"`
	OwnerID     primitive.ObjectID   `bson:"owner_id,omitempty" json:"owner_id"`
	TeamIDs     []primitive.ObjectID `bson:"team,omitempty" json:"teams_id"`
	Status      string               `bson:"status" json:"status"`
}
