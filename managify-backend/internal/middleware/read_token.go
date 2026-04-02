package middleware

import (
	"fmt"
	"managify/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExtractUserFromClaims converts JWT claims into a User model.
func ExtractUserFromClaims(claims jwt.MapClaims) (*models.User, error) {
	id, err := claimToObjectID(claims, "id")
	if err != nil {
		return nil, err
	}

	name, err := claimToString(claims, "name")
	if err != nil {
		return nil, err
	}

	email, err := claimToString(claims, "email")
	if err != nil {
		return nil, err
	}

	isAdmin := claimToBool(claims, "is_admin", false)

	return &models.User{
		ID:       id,
		FullName: name,
		Email:    email,
		IsAdmin:  isAdmin,
	}, nil
}

// Helpers
func claimToString(claims jwt.MapClaims, key string) (string, error) {
	v, ok := claims[key].(string)
	if !ok || v == "" {
		return "", fmt.Errorf("%s claim missing or invalid", key)
	}
	return v, nil
}

func claimToBool(claims jwt.MapClaims, key string, defaultVal bool) bool {
	v, ok := claims[key].(bool)
	if !ok {
		return defaultVal
	}
	return v
}

func claimToObjectID(claims jwt.MapClaims, key string) (primitive.ObjectID, error) {
	idStr, err := claimToString(claims, key)
	if err != nil {
		return primitive.NilObjectID, err
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("invalid ObjectID for %s: %v", key, err)
	}
	return id, nil
}
