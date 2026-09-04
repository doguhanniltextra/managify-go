package router

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRouter_RESTfulRouteRegistration(t *testing.T) {
	app := fiber.New()

	// Register all routers
	RouterUser(app)
	RouterAdmin(app)
	RouterProject(app)
	RouterInvite(app)
	RouterRole(app)
	RouterStatus(app)
	RouterIssue(app)
	RouterLogger(app)

	routes := app.GetRoutes()

	type expectedRoute struct {
		method string
		path   string
	}

	expected := []expectedRoute{
		// Project REST endpoints
		{method: "POST", path: "/v1/projects"},
		{method: "GET", path: "/v1/projects/:id"},
		{method: "DELETE", path: "/v1/projects/:id"},
		{method: "DELETE", path: "/v1/projects/member/:memberId"},

		// Admin REST endpoints
		{method: "GET", path: "/v1/admin/users"},
		{method: "GET", path: "/v1/admin/users/:id"},
		{method: "DELETE", path: "/v1/admin/users/:id"},
		{method: "GET", path: "/v1/admin/projects"},
		{method: "GET", path: "/v1/admin/roles"},

		// Role REST endpoints
		{method: "POST", path: "/v1/roles"},
		{method: "DELETE", path: "/v1/roles/:id/:project"},

		// Status REST endpoints
		{method: "POST", path: "/v1/statuses"},
		{method: "DELETE", path: "/v1/statuses/:id/:project"},

		// Issue REST endpoints
		{method: "POST", path: "/v1/issues"},
		{method: "DELETE", path: "/v1/issues/:id"},
		{method: "GET", path: "/v1/issues/status/:statusID"},
		{method: "PUT", path: "/v1/issues/:issueID/status/:statusID"},

		// Invite REST endpoints
		{method: "POST", path: "/v1/invites"},
		{method: "GET", path: "/v1/invites/:id"},
		{method: "PUT", path: "/v1/invites/:inviteId/respond"},
	}

	for _, exp := range expected {
		found := false
		for _, r := range routes {
			if r.Method == exp.method && (r.Path == exp.path || r.Path == exp.path+"/") {
				found = true
				break
			}
		}
		assert.Truef(t, found, "Expected route %s %s not registered in Fiber", exp.method, exp.path)
	}
}
