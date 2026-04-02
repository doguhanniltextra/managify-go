package handler

import (
	"managify/internal/service"

	"github.com/gofiber/fiber/v2"
)

func GetLogsHandler(c *fiber.Ctx) error {
	projectID := c.Params("projectId")
	if projectID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Project ID is required"})
	}

	logs, err := service.GetLogService().GetLogsByProjectID(projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch logs"})
	}

	return c.JSON(fiber.Map{"logs": logs})
}

func GetLogsHandlerByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")
	if userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User ID is required"})
	}

	logs, err := service.GetLogService().GetLogsByUserId(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch logs"})
	}

	return c.JSON(fiber.Map{"logs": logs})
}
