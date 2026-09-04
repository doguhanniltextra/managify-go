package handler

import (
	"errors"
	"managify/constant"
	"managify/internal/domain"

	"github.com/gofiber/fiber/v2"
)

// HandleServiceError maps domain errors to proper HTTP responses.
func HandleServiceError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constant.ErrNotFound,
			"error":   err.Error(),
		})
	case errors.Is(err, domain.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": constant.ErrForbidden,
			"error":   err.Error(),
		})
	case errors.Is(err, domain.ErrBadRequest):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   err.Error(),
		})
	case errors.Is(err, domain.ErrConflict):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": constant.ErrConflict,
			"error":   err.Error(),
		})
	case errors.Is(err, domain.ErrUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
			"error":   err.Error(),
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
			"error":   err.Error(),
		})
	}
}
