package handler

import (
	"managify/constant"
	"managify/internal/service"
	"managify/models"
	"managify/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create a new status
// @Description Creates a new status in the system.
// @Tags Statuses
// @Accept json
// @Produce json
// @Param status body models.Status true "Status to create"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /statuses [post]
func CreateStatusHandler(c *fiber.Ctx) error {
	var status models.Status

	if err := c.BodyParser(&status); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   err.Error(),
		})
	}

	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	status.CreatorID = user.ID

	res, err := service.GetStatusService().CreateStatus(&status)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessCreated,
		"data":    res,
	})
}

// @Summary Delete a status
// @Description Deletes a status by its ID for a specific project.
// @Tags Statuses
// @Produce json
// @Param id path string true "Status ID"
// @Param project path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /statuses/{id}/project/{project} [delete]
func DeleteStatusHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	projectParam := c.Params("project")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	statusID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}
	projectID, err := primitive.ObjectIDFromHex(projectParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}
	user, ok := utils.GetUserLocal(c) // User
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	err = service.GetStatusService().DeleteStatus(statusID, projectID, user.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessDeleted,
	})
}
