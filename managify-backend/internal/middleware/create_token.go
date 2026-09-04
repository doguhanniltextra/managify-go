package middleware

import (
	"managify/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func InitJWTMiddleware(key string) {
	secretKey = []byte(key)
}

func CreateToken(user *models.User) (string, error) {
	mapClaims := jwt.MapClaims{
		"id":       user.ID,
		"name":     user.FullName,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"iss":      "user",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
