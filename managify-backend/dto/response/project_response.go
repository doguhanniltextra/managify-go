package response

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusWithIssues struct {
	ID        primitive.ObjectID   `json:"id"`
	ProjectID primitive.ObjectID   `json:"project_id"`
	Name      string               `json:"name"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at,omitempty"`
	IssuesID  []primitive.ObjectID `bson:"issues" json:"issues_id"`
}
