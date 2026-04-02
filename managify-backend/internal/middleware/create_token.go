package middleware

import (
	"managify/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func CreateToken(user *models.User) (string, error) {

	var log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
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

	log.Info(claims)
	log.Info(claims.Header)
	log.Info(claims.Signature)

	tokenString, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
