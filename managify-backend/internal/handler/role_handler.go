package handler

import (
	"fmt"

	"managify/constant"
	"managify/internal/service"
	"managify/models"
	"managify/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create a new role
// @Description Creates a new role in a project. Only project owners can add roles.
// @Tags Projects
// @Accept json
// @Produce json
// @Param role body models.Role true "Role to create"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /roles [post]
func CreateRoleHandler(c *fiber.Ctx) error {
	var role models.Role

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	roleUserID := role.UserID
	roleProjectID := role.ProjectID
	RoleName := role.RoleName

	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	fmt.Println(user.ID)

	isOwner, err := service.GetProjectService().IsOwner(user.ID, roleProjectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}
	if !isOwner {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": constant.ErrForbidden,
			"error":   "You are not owner",
		})
	}

	res, err := service.GetRoleService().AddRole(roleUserID, roleProjectID, RoleName)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessCreated,
		"data":    res,
	})
}

// @Summary Delete a role
// @Description Deletes a role by its ID.
// @Tags Projects
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /roles/{id} [delete]
func DeleteRoleHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   "missing id parameter",
		})
	}

	roleID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   "invalid id format",
		})
	}

	err = service.GetRoleService().DeleteRole(roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessDeleted,
	})
}
