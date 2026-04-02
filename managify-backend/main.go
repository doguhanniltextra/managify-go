package main

import (
	"fmt"
	"managify/constant"
	"managify/database"
	"managify/internal/middleware"
	"managify/internal/router"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/pprof"

	_ "managify/internal/handler"
	_ "managify/swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
)

// @Title Managify API
// @Version 1.0
// @Description This is the API documentation for the Managify project management application. It provides endpoints for managing projects, issues, users, roles, statuses, and project invites.
// @Host localhost:3000
// @Contact.name Doguhan İlter
// @Contact.email doguhannilt@gmail.com
// @BasePath /
func main() {
	middleware.InitLogger()
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, relying on environment variables")
	}

	if err := database.Connect(); err != nil {
		logrus.Infoln("Database connection failed: ", err)
	}

	app := fiber.New()

	host := os.Getenv("VUE_HOST")
	var allowOrigins string
	var allowCredentials bool

	if host == "" {
		logrus.Warn("VUE_HOST is not set, defaulting to wildcard origins without credentials")
		allowOrigins = "*"
		allowCredentials = false
	} else {
		allowOrigins = host
		allowCredentials = true
	}

	allowMethods := "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS"
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: allowCredentials,
	}))

	apiLimiter(app)

	app.Use(middleware.MetricMiddleware)

	// pprof for performance monitoring
	app.Use(pprof.New())
	router.Routers(app)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	logrus.Infof("Starting server on %s", addr)
	if err := app.Listen(addr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}

}

func apiLimiter(app *fiber.App) {

	apiMaxLimiter, _ := strconv.Atoi(os.Getenv("API_MAX_LIMITER"))
	expirationSec, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_EXPIRATION"))

	app.Use(limiter.New(limiter.Config{
		Max:        apiMaxLimiter,
		Expiration: time.Duration(expirationSec) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": constant.ErrTooManyRequests,
			})
		},
	}))
}
