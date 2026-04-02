package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	FullName          string               `bson:"full_name" json:"full_name"`
	Email             string               `bson:"email" json:"email"`
	Password          string               `bson:"password" json:"password"`
	GoogleID          string               `bson:"google_id,omitempty" json:"-"`
	AuthProvider      string               `bson:"auth_provider,omitempty" json:"auth_provider"`
	AssignedIssues    []primitive.ObjectID `bson:"assigned_issues,omitempty" json:"-"`
	ProjectSize       int                  `bson:"project_size" json:"project_size"`
	Subscriptions     []primitive.ObjectID `bson:"subscriptions,omitempty" json:"-"`
	OwnedProjects     []primitive.ObjectID `bson:"owned_projects,omitempty" json:"-"`
	TeamProjects      []primitive.ObjectID `bson:"team_projects,omitempty" json:"-"`
	IsAdmin           bool                 `bson:"is_admin" json:"is_admin"`
	VerificationToken string               `bson:"verificationtoken,omitempty" json:"verificationtoken,omitempty"`
	IsVerified        bool                 `bson:"isverified" json:"isverified"`
}
