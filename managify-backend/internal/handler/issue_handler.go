package handler

import (
	"managify/constant"
	"managify/internal/service"
	"managify/models"
	"managify/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Get issues by status ID
// @Description Retrieves all issues associated with a specific status ID.
// @Tags Issues
// @Produce json
// @Param statusID path string true "Status ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issues/status/{statusID} [get]
func GetIssuesByStatusHandler(c *fiber.Ctx) error {
	statusIDHex := c.Params("statusID")
	statusID, err := primitive.ObjectIDFromHex(statusIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": constant.ErrBadRequest})
	}

	issues, err := service.GetIssueService().GetIssuesByStatusID(statusID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": constant.ErrInternalServer})
	}

	data := make([]fiber.Map, 0, len(issues))
	for _, i := range issues {
		data = append(data, fiber.Map{
			"id":          i.ID,
			"title":       i.Title,
			"description": i.Description,
			"priority":    i.Priority,
			"due_date":    i.DueDate,
			"status_id":   i.StatusID,
			"project_id":  i.ProjectID,
		})
	}

	return c.JSON(fiber.Map{
		"message": constant.SuccessFetched,
		"data":    data,
	})
}

// @Summary Create a new issue
// @Description Creates a new issue in the system.
// @Tags Issues
// @Accept json
// @Produce json
// @Param issue body models.Issue true "Issue to create"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issues [post]
func CreateIssueHandler(c *fiber.Ctx) error {

	var issue models.Issue

	if err := c.BodyParser(&issue); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	user, ok := utils.GetUserLocal(c)
	if !ok {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	res, err := service.GetIssueService().CreateIssue(&issue, user.ID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessCreated,
		"data":    res,
	})
}

// @Summary Delete an issue
// @Description Deletes an issue by its ID.
// @Tags Issues
// @Produce json
// @Param id path string true "Issue ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issues/{id} [delete]
func DeleteIssueHandler(c *fiber.Ctx) error {
	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	idStr := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	err = service.GetIssueService().DeleteIssue(objID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": constant.ErrUnauthorized,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessDeleted,
	})
}

// @Summary Update issue status
// @Description Updates the status of an issue.
// @Tags Issues
// @Produce json
// @Param issueID path string true "Issue ID"
// @Param statusID path string true "New Status ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issues/{issueID}/status/{statusID} [patch]
func UpdateIssueStatusHandler(c *fiber.Ctx) error {

	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	issueIDHex := c.Params("issueID")
	newStatusIDHex := c.Params("statusID")

	issueID, err := primitive.ObjectIDFromHex(issueIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	newStatusID, err := primitive.ObjectIDFromHex(newStatusIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": constant.ErrBadRequest,
		})
	}

	updatedIssue, err := service.GetIssueService().UpdateIssueStatus(issueID, newStatusID, user.ID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessUpdated,
		"data": fiber.Map{
			"id":          updatedIssue.ID,
			"title":       updatedIssue.Title,
			"description": updatedIssue.Description,
			"priority":    updatedIssue.Priority,
			"due_date":    updatedIssue.DueDate,
			"status_id":   updatedIssue.StatusID.Hex(),
			"project_id":  updatedIssue.ProjectID,
		},
	})
}

// @Summary Get oncoming issues for a project
// @Description Retrieves issues that are due soon for a specific project.
// @Tags Issues
// @Produce json
// @Param projectID path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /issues/project/{projectID}/oncoming [get]
func GetOncomingIssuesHandler(c *fiber.Ctx) error {
	projectIDHex := c.Params("projectID")
	projectID, err := primitive.ObjectIDFromHex(projectIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": constant.ErrBadRequest})
	}

	issues, err := service.GetIssueService().GetOncomingIssues(projectID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": constant.ErrInternalServer})
	}

	issueResponse := make([]fiber.Map, 0, len(issues))
	for _, issues := range issues {
		issueResponse = append(issueResponse, fiber.Map{
			"id":          issues.ID,
			"title":       issues.Title,
			"description": issues.Description,
			"due_date":    issues.DueDate,
		})
	}
	return c.JSON(fiber.Map{
		"message": constant.SuccessFetched,
		"data":    issueResponse,
	})
}
