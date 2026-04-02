package handler

import (
	"managify/constant"
	"managify/internal/service"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get Google OAuth URL
// @Description Returns the Google OAuth consent screen URL for the frontend to redirect to.
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/google [get]
func GoogleAuthURLHandler(c *fiber.Ctx) error {
	url := service.GetGoogleAuthService().GetGoogleAuthURL()
	return c.JSON(fiber.Map{
		"url": url,
	})
}

// @Summary Google OAuth Callback
// @Description Receives the authorization code from Google, exchanges it for tokens, and returns a Managify JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body object true "Google authorization code" example({"code": "4/0ABC123..."})
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/google/callback [post]
func GoogleCallbackHandler(c *fiber.Ctx) error {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   err.Error(),
		})
	}

	if req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   "authorization code is required",
		})
	}

	user, token, err := service.GetGoogleAuthService().HandleCallback(req.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessOperation,
		"token":   token,
		"email":   user.Email,
		"name":    user.FullName,
	})
}
