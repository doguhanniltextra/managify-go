package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PriorityType string
type StatusType string

const (
	Default  PriorityType = "DEFAULT"
	Medium   PriorityType = "MEDIUM"
	High     PriorityType = "HIGH"
	Urgent   PriorityType = "URGENT"
	Critical PriorityType = "CRITICAL"
)

const (
	TODO        StatusType = "TODO"
	IN_PROGRESS StatusType = "IN_PROGRESS"
	REVIEW      StatusType = "REVIEW"
	DONE        StatusType = "DONE"
	BLOCKED     StatusType = "BLOCKED"
)

type Issue struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	Status      StatusType           `bson:"status" json:"status"`
	ProjectID   primitive.ObjectID   `bson:"project_id,omitempty" json:"project_id"`
	Priority    PriorityType         `bson:"priority" json:"priority"`
	DueDate     string               `bson:"due_date,omitempty" json:"due_date"`
	Tags        []string             `bson:"tags,omitempty" json:"tags"`
	StatusID    primitive.ObjectID   `bson:"status_id,omitempty" json:"status_id"`
	CommentIDs  []primitive.ObjectID `bson:"comments,omitempty" json:"-"`
}
