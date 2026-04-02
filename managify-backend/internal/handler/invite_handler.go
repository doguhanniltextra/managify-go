package handler

import (
	"managify/constant"
	"managify/dto/request"
	"managify/internal/service"
	"managify/models"
	"managify/utils"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create a new project invite
// @Description Creates a new invite for a project.
// @Tags ProjectInvites
// @Accept json
// @Produce json
// @Param invite body request.ProjectInviteRequest true "Project Invite Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /project-invites [post]
func CreateProjectInviteHandler(c *fiber.Ctx) error {
	var req request.ProjectInviteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	invite, err := service.CreateProjectInvite(user.ID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	return c.JSON(fiber.Map{
		"message": constant.SuccessCreated,
		"invite":  invite,
	})
}

// @Summary Respond to a project invite
// @Description Accepts or declines a project invite.
// @Tags ProjectInvites
// @Produce json
// @Param inviteId path string true "Invite ID"
// @Param action query string true "Action (accept/decline)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /project-invites/{inviteId}/respond [patch]
func RespondProjectInviteHandler(c *fiber.Ctx) error {
	inviteIDHex := c.Params("inviteId")
	inviteID, err := primitive.ObjectIDFromHex(inviteIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	userVal := c.Locals("user")
	user, ok := userVal.(*models.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": constant.ErrUnauthorized})
	}

	action := c.Query("action")
	accept := false
	if action == "accept" {
		accept = true
	} else if action != "decline" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": constant.ErrBadRequest})
	}

	invite, err := service.RespondProjectInvite(user.ID, inviteID, accept)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": constant.ErrBadRequest})
	}

	return c.JSON(fiber.Map{
		"message": constant.SuccessUpdated,
		"invite":  invite,
	})
}

// @Summary Get project invites by ID
// @Description Retrieves project invites by user or project ID.
// @Tags ProjectInvites
// @Produce json
// @Param id path string true "User or Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /project-invites/{id} [get]
func GetInviteHandlerById(c *fiber.Ctx) error {
	rcvStr := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(rcvStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	models, err := service.GetProjectInvites(objID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessFetched,
		"data":    models,
	})
}
