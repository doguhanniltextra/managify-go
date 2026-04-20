package router

import (
	"managify/internal/handler"
	"managify/internal/middleware"
	"managify/internal/router/routes"
	"os"

	"managify/internal/validation"
	"managify/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Routers(app *fiber.App) {
	RouterGoogleAuth(app)
	RouterUser(app)
	RouterAdmin(app)
	RouterProject(app)
	RouterInvite(app)
	RouterRole(app)
	RouterIssue(app)
	RouterStatus(app)
	RouterLogger(app)
	RouterSwagger(app)
	RouterMetrics(app)

	// Health check for Load Balancer (Deep Health Check)
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := database.CheckHealth(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "error",
				"error":  "database connection lost",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})
}

func RouterUser(app *fiber.App) {
	api := app.Group(routes.UserBase)
	api.Get(routes.UserVerifyEmail, handler.VerifyEmailHandler)
	api.Post(routes.UserRegister, validation.CreateRegisterValidator, handler.CreateRegisterHandler)
	api.Post(routes.UserAuth, validation.AuthValidator, handler.LoginHandler)
	api.Get(routes.UserGetById, middleware.AuthMiddleware, handler.GetUserByIdHandler)

}

func RouterAdmin(app *fiber.App) {
	api := app.Group(routes.AdminBase, middleware.AuthMiddleware, middleware.AdminMiddleware)

	api.Get(routes.AdminGetUsers, handler.GetUsersHandler)
	api.Get(routes.AdminGetUser, handler.GetUserById)
	api.Get(routes.AdminGetProjects, handler.GetProjectsHandler)
	api.Get(routes.AdminGetRoles, handler.GetRolesHandler)
	api.Delete(routes.AdminDelete, handler.DeleteUserById)
}

func RouterProject(app *fiber.App) {
	api := app.Group(routes.ProjectBase, middleware.AuthMiddleware)

	api.Post(routes.ProjectCreate, handler.CreateProjectHandler)
	api.Delete(routes.ProjectDelete, handler.DeleteProjectHandler)
	api.Get(routes.ProjectGet, handler.GetProjectHandler)
	api.Delete(routes.ProjectMemberDelete, handler.DeleteMemberFromProjectByIdHandler)
}

func RouterInvite(app *fiber.App) {
	api := app.Group(routes.InviteBase, middleware.AuthMiddleware)

	api.Get(routes.InviteGetById, handler.GetInviteHandlerById)
	api.Post(routes.InviteCreate, handler.CreateProjectInviteHandler)
	api.Put(routes.InviteRespond, handler.RespondProjectInviteHandler)
}

func RouterRole(app *fiber.App) {
	api := app.Group(routes.RoleBase, middleware.AuthMiddleware)

	api.Post(routes.RoleCreate, handler.CreateRoleHandler)
	api.Delete(routes.RoleDelete, handler.DeleteRoleHandler)
}

func RouterStatus(app *fiber.App) {
	api := app.Group(routes.StatusBase, middleware.AuthMiddleware)

	api.Post(routes.StatusCreate, validation.CreateStatusValidator, handler.CreateStatusHandler)
	api.Delete(routes.StatusDelete, handler.DeleteStatusHandler)
}

func RouterIssue(app *fiber.App) {
	api := app.Group(routes.IssueBase, middleware.AuthMiddleware)

	api.Post(routes.IssueCreate, handler.CreateIssueHandler)
	api.Delete(routes.IssueDelete, handler.DeleteIssueHandler)
	api.Get(routes.IssuesGet, handler.GetIssuesByStatusHandler)
	api.Put(routes.IssueUpdate, handler.UpdateIssueStatusHandler)
	api.Get(routes.IssueGetOnDue, handler.GetOncomingIssuesHandler)
}

func RouterLogger(app *fiber.App) {
	api := app.Group(routes.LoggerBase, middleware.AuthMiddleware)

	api.Get(routes.LoggerGet, handler.GetLogsHandlerByUserId)
}

func RouterSwagger(app *fiber.App) {
	api := app.Group(routes.SwaggerBase)

	if os.Getenv("SWAGGER") == "true" {
		api.Get("/*", swagger.HandlerDefault)
	}
}

func RouterMetrics(app *fiber.App) {
	api := app.Group(routes.MetricsBase)

	if os.Getenv("METRICS") == "true" {
		api.Get("/", handler.MetricsHandler)
	}
}

func RouterGoogleAuth(app *fiber.App) {
	api := app.Group(routes.GoogleAuthBase)

	api.Get(routes.GoogleAuthURL, handler.GoogleAuthURLHandler)
	api.Post(routes.GoogleAuthCallback, handler.GoogleCallbackHandler)
}
