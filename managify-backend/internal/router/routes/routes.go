package routes

const (
	version = "/v1"

	// Google OAuth endpoints
	GoogleAuthBase     = version + "/auth"
	GoogleAuthURL      = "/google"
	GoogleAuthCallback = "/google/callback"

	// User endpoints
	UserBase        = version + "/users"
	UserRegister    = "/register"
	UserAuth        = "/auth"
	UserGetById     = "/:id"
	UserVerifyEmail = "/verify-email"

	// Admin endpoints
	AdminBase        = version + "/admin"
	AdminUsers       = "/users"
	AdminUser        = "/users/:id"
	AdminProjects    = "/projects"
	AdminRoles       = "/roles"

	// Project endpoints
	ProjectBase         = version + "/projects"
	ProjectCreate       = "/"
	ProjectGet          = "/:id"
	ProjectDelete       = "/:id"
	ProjectMemberDelete = "/member/:memberId"

	// Project invite endpoints
	InviteBase    = version + "/invites"
	InviteCreate  = "/"
	InviteGetById = "/:id"
	InviteRespond = "/:inviteId/respond"

	// Project role endpoints
	RoleBase   = version + "/roles"
	RoleCreate = "/"
	RoleDelete = "/:id/:project"

	// Project status endpoints
	StatusBase   = version + "/statuses"
	StatusCreate = "/"
	StatusDelete = "/:id/:project"

	// Project issue endpoints
	IssueBase     = version + "/issues"
	IssueCreate   = "/"
	IssueDelete   = "/:id"
	IssuesGet     = "/status/:statusID"
	IssueUpdate   = "/:issueID/status/:statusID"
	IssueGetOnDue = "/due-today/:projectID"

	// Log endpoint
	LoggerBase = version + "/logger"
	LoggerGet  = "/:userId"

	// Swagger endpoint
	SwaggerBase = version + "/swagger"

	// Metrics endpoint
	MetricsBase = version + "/metrics"
)
