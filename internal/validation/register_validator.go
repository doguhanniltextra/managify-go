package validation

import (
	"context"
	"managify/database"
	"managify/internal/repository"
	"managify/models"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func CreateRegisterValidator(c *fiber.Ctx) error {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.InfoLevel)

	var user models.User

	// Body parse
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Password validation
	if user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password is required",
		})
	}

	if len(user.Password) < 6 || len(user.Password) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must be between 6 and 20 characters",
		})
	}

	if CheckPasswordComplexity(user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password must contain at least 1 number, 1 uppercase letter, and 1 special character",
		})
	}

	// Email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid email format",
		})
	}

	// DB uniqueness checks
	userRepo := repository.NewUserRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingUserByEmail, err := userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}
	if existingUserByEmail != nil {
		log.Warnf("Email already exists: %s", user.Email)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	existingUserByFullName, err := userRepo.FindByFullName(ctx, user.FullName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}
	if existingUserByFullName != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Name already exists",
		})
	}

	return c.Next()
}

func CheckPasswordComplexity(password string) bool {
	hasNumber := false
	hasUpper := false
	hasSpecial := false

	for _, c := range password {
		switch {
		case c >= '0' && c <= '9':
			hasNumber = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case (c >= '!' && c <= '/') ||
			(c >= ':' && c <= '@') ||
			(c >= '[' && c <= '`') ||
			(c >= '{' && c <= '~'):
			hasSpecial = true
		}
	}

	return hasNumber && hasUpper && hasSpecial
}
