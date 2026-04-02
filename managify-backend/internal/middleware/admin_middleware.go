package middleware

import (
	"managify/constant"
	"managify/models"

	"github.com/gofiber/fiber/v2"
)

// AdminMiddleware ensures only admins can access the route
func AdminMiddleware(c *fiber.Ctx) error {
	userVal := c.Locals("user")
	if userVal == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	user, ok := userVal.(*models.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	if !user.IsAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": constant.ErrForbidden,
		})
	}

	return c.Next()
}
