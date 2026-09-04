package middleware

import (
	"managify/models"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateToken(t *testing.T) {
	testSecret := "super-secret-jwt-key-for-unit-testing"
	InitJWTMiddleware(testSecret)

	user := &models.User{
		ID:       primitive.NewObjectID(),
		FullName: "Test User",
		Email:    "test@example.com",
		IsAdmin:  true,
	}

	tokenString, err := CreateToken(user)
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Parse token to verify signature and claims
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		assert.Equal(t, jwt.SigningMethodHS256, token.Method)
		return []byte(testSecret), nil
	})

	require.NoError(t, err)
	require.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, user.FullName, claims["name"])
	assert.Equal(t, user.Email, claims["email"])
	assert.Equal(t, user.IsAdmin, claims["is_admin"])
	assert.Equal(t, "user", claims["iss"])

	// Verify expiration and issued-at times
	exp, ok := claims["exp"].(float64)
	require.True(t, ok)
	assert.True(t, int64(exp) > time.Now().Unix())

	iat, ok := claims["iat"].(float64)
	require.True(t, ok)
	assert.True(t, int64(iat) <= time.Now().Unix())
}
