package middleware

import (
	"fmt"
	"managify/constant"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	// Bearer kısmını ayır
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	// Token parse
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	user, err := ExtractUserFromClaims(claims)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	c.Locals("user", user)

	return c.Next()
}
