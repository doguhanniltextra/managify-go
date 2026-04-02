package validation

import (
	"context"
	"fmt"
	"managify/database"
	"managify/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProjectValidator(c *fiber.Ctx) error {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.InfoLevel)

	var project models.Project

	// Body parse
	if err := c.BodyParser(&project); err != nil {
		log.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	log.Debugf("Parsed request body: %+v", project)

	if err := validateUserId(project.OwnerID); err != nil {
		log.WithError(err).Warnf("Invalid owner id: %s", project.OwnerID.Hex())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Name validation
	if project.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Project name is required",
		})
	}
	if len(project.Name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Project name must be at most 100 characters",
		})
	}

	if len(project.Tags) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "At least one tag is required",
		})
	}

	if project.Description == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Project description is required",
		})
	}
	if len(project.Description) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Project description must be at most 500 characters",
		})
	}

	if project.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Project category is required",
		})
	}

	return c.Next()
}

func validateUserId(id primitive.ObjectID) error {
	collection := database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to check user id: %v", err)
	}

	if count == 0 {
		return fmt.Errorf("user with id %s does not exist", id.Hex())
	}

	return nil
}
