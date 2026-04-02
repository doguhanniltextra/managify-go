package validation

import (
	"managify/models"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func CreateStatusValidator(c *fiber.Ctx) error {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.InfoLevel)

	var status models.Status

	if err := c.BodyParser(&status); err != nil {
		log.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if status.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status name is required",
		})
	}
	if len(status.Name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Status name must be at most 100 characters",
		})
	}

	return c.Next()
}
