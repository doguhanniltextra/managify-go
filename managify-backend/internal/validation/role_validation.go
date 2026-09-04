package validation

import (
	"managify/models"

	"github.com/gofiber/fiber/v2"
)

func CreateRoleValidator(c *fiber.Ctx) error {

	var role models.Role

	if err := c.BodyParser(&role); err != nil {
		log.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if role.RoleName == "" {
		log.Error("Role name is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Role name is required",
		})
	}

	return c.Next()
}
