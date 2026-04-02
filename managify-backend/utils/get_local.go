package utils

import (
	"managify/models"

	"github.com/gofiber/fiber/v2"
)

func GetUserLocal(c *fiber.Ctx) (*models.User, bool) {
	userVal := c.Locals("user")
	user, ok := userVal.(*models.User)
	if !ok {
		return nil, false
	}
	return user, true
}
