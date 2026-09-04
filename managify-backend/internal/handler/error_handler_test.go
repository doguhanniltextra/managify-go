package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"managify/constant"
	"managify/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleServiceError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Not Found Error",
			err:            domain.ErrNotFound,
			expectedStatus: fiber.StatusNotFound,
			expectedMsg:    constant.ErrNotFound,
		},
		{
			name:           "Wrapped Not Found Error",
			err:            fmt.Errorf("%w: project not found", domain.ErrNotFound),
			expectedStatus: fiber.StatusNotFound,
			expectedMsg:    constant.ErrNotFound,
		},
		{
			name:           "Forbidden Error",
			err:            fmt.Errorf("%w: user is not in project", domain.ErrForbidden),
			expectedStatus: fiber.StatusForbidden,
			expectedMsg:    constant.ErrForbidden,
		},
		{
			name:           "Bad Request Error",
			err:            fmt.Errorf("%w: invalid id format", domain.ErrBadRequest),
			expectedStatus: fiber.StatusBadRequest,
			expectedMsg:    constant.ErrBadRequest,
		},
		{
			name:           "Conflict Error",
			err:            fmt.Errorf("%w: already exists", domain.ErrConflict),
			expectedStatus: fiber.StatusConflict,
			expectedMsg:    constant.ErrConflict,
		},
		{
			name:           "Unauthorized Error",
			err:            fmt.Errorf("%w: invalid token", domain.ErrUnauthorized),
			expectedStatus: fiber.StatusUnauthorized,
			expectedMsg:    constant.ErrUnauthorized,
		},
		{
			name:           "Internal Unknown Error",
			err:            errors.New("unexpected database crash"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedMsg:    constant.ErrInternalServer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				return HandleServiceError(c, tc.err)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var res map[string]interface{}
			err = json.Unmarshal(body, &res)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedMsg, res["message"])
			assert.NotEmpty(t, res["error"])
		})
	}
}

func TestHandleServiceError_Nil(t *testing.T) {
	app := fiber.New()
	app.Get("/test-nil", func(c *fiber.Ctx) error {
		err := HandleServiceError(c, nil)
		if err == nil {
			return c.SendString("ok")
		}
		return err
	})

	req := httptest.NewRequest("GET", "/test-nil", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
