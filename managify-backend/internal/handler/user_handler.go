package handler

import (
	"managify/constant"
	"managify/dto/request"
	"managify/internal/service"
	"managify/models"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// @Summary Register a new user
// @Description Registers a new user in the system.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register [post]
func CreateRegisterHandler(c *fiber.Ctx) error {
	var user models.User

	// Body parse
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   err.Error(),
		})
	}

	createdUser, token, err := service.GetUserService().CreateUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
			"error":   err.Error(),
		})
	}

	subscriptionStartDate := time.Now()
	subscriptionEndDate := time.Now()
	planType := models.PlanBasic
	isValid := true

	subscriptionMethod := models.Subscription{

		SubscriptionStartDate: subscriptionStartDate,
		SubscriptionEndDate:   subscriptionEndDate,
		PlanType:              planType,
		IsValid:               isValid,
		UserID:                user.ID,
	}

	var (
		wg              sync.WaitGroup
		subscription    *models.Subscription
		subscriptionErr error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		subscription, err = service.GetSubscriptionService().CreateSubscription(&subscriptionMethod)
		if err != nil {
			subscriptionErr = err
			return
		}
	}()

	wg.Wait()

	if subscriptionErr != nil {
		return subscriptionErr
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":      constant.SuccessCreated,
		"token":        token,
		"userEmail":    createdUser.Email,
		"subscription": subscription.PlanType,
	})
}

// @Summary User login
// @Description Authenticates a user and returns a JWT token.
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body request.UserLoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func LoginHandler(c *fiber.Ctx) error {
	var req request.UserLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
			"error":   err.Error(),
		})
	}

	res, err := service.GetUserService().Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessOperation,
		"email":   res.Email,
		"name":    res.FullName,
		"token":   res.Token,
	})
}

// @Summary Get user by ID
// @Description Retrieves a user by their ID, along with associated projects and subscription details.
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [get]
func GetUserByIdHandler(c *fiber.Ctx) error {
	userIDHex := c.Params("id")

	var (
		wg sync.WaitGroup

		user       *models.User
		project    any
		sub        any
		userErr    error
		projectErr error
		subErr     error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		user, userErr = service.GetUserService().GetUserById(userIDHex)
	}()

	go func() {
		defer wg.Done()
		project, projectErr = service.GetProjectService().GetProjectsByUserId(userIDHex)
	}()

	go func() {
		defer wg.Done()
		sub, subErr = service.GetSubscriptionService().GetByUserId(userIDHex)
	}()

	wg.Wait()

	if userErr != nil || projectErr != nil || subErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": constant.ErrNotFound,
		})
	}

	data := fiber.Map{
		"user":            user,
		"isVerified":      user.IsVerified,
		"validationToken": user.VerificationToken,
		"project":         project,
		"subscription":    sub,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessOperation,
		"data":    data,
	})
}

// @Summary Verify user email
// @Description Verifies a user's email using a token.
// @Tags Users
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /users/verify-email [get]
func VerifyEmailHandler(c *fiber.Ctx) error {
	token := c.Query("token")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Token missing"})
	}

	user, err := service.GetUserService().VerifyEmail(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid or expired token"})
	}

	return c.JSON(fiber.Map{"message": "Email verified", "user": user.Email})
}
