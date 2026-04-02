package routes

const (
	version = "/v1"

	// Google OAuth endpoints
	GoogleAuthBase     = version + "/auth"
	GoogleAuthURL      = "/google"
	GoogleAuthCallback = "/google/callback"

	// User endpoints
	UserBase     = version + "/users"
	UserRegister = "/register"
	UserAuth     = "/auth"
	UserGetById  = "/:id"

	UserVerifyEmail = "/verify-email"

	// Admin endpoints
	AdminBase        = version + "/admin"
	AdminGetUsers    = "/get-users"
	AdminGetUser     = "/get-user/:id"
	AdminDelete      = "/delete-user/:id"
	AdminGetProjects = "/get-projects"
	AdminGetRoles    = "/get-roles"

	// Project endpoints

	ProjectBase         = version + "/project"
	ProjectCreate       = "/create-project"
	ProjectDelete       = "/delete-project/:id"
	ProjectGet          = "/projects/:id"
	ProjectMemberDelete = "/projects/member/:memberId"

	// Project invite endpoints

	InviteBase    = version + "/invite"
	InviteCreate  = "/project-invite"
	InviteRespond = "/project-invite/:inviteId/respond"
	InviteGetById = "/project-invite/:id"

	// Project role endpoints

	RoleBase   = version + "/role"
	RoleCreate = "/create-role"
	RoleDelete = "/delete-role/:id/:project"

	// Project status endpoints

	StatusBase   = version + "/status"
	StatusCreate = "/create-status"
	StatusDelete = "/delete-status/:id/:project"

	// Project issue endpoints

	IssueBase     = version + "/issue"
	IssueCreate   = "/create-issue"
	IssueDelete   = "/delete-issue/:id"
	IssuesGet     = "/get/:statusID"
	IssueUpdate   = "/update-status/:issueID/:statusID"
	IssueGetOnDue = "/due-today/:projectID"

	// Log endpoint
	LoggerBase = version + "/logger"
	LoggerGet  = "/:userId"

	// Swagger endpoint
	SwaggerBase = version + "/swagger"

	// Metrics endpoint
	MetricsBase = version + "/metrics"
)
