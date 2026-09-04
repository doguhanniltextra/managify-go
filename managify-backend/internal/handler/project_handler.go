package handler

import (
	"fmt"
	"managify/constant"
	"managify/internal/service"
	"managify/models"
	"managify/utils"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"managify/dto/response"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create a new project
// @Description Creates a new project in the system.
// @Tags Projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project to create"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects [post]
func CreateProjectHandler(c *fiber.Ctx) error {
	var project models.Project

	if err := c.BodyParser(&project); err != nil {
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

	res, err := service.GetProjectService().CreateProject(c.UserContext(), &project, user)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessCreated,
		"data":    res,
	})
}

// @Summary Delete a project
// @Description Deletes a project by its ID.
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects/{id} [delete]
func DeleteProjectHandler(c *fiber.Ctx) error {

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

	err = service.GetProjectService().DeleteProjectById(c.UserContext(), objID, user)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessDeleted,
	})
}


// @Summary Get a project by ID
// @Description Retrieves a project by its ID, including its statuses and team members.
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects/{id} [get]
func GetProjectHandler(c *fiber.Ctx) error {
	start := time.Now()
	projectIDHex := c.Params("id")
	projectID, err := primitive.ObjectIDFromHex(projectIDHex)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": constant.ErrBadRequest})
	}

	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": constant.ErrInternalServer})
	}

	project, err := service.GetProjectService().GetProject(c.UserContext(), projectID, user)
	if err != nil {
		return HandleServiceError(c, err)
	}

	var (
		statuses    []*models.Status
		teamMembers []models.User
		statusErr   error
		teamErr     error
		wg          sync.WaitGroup
	)

	wg.Add(2)

	// Statuses parallel fetch
	go func() {
		defer wg.Done()
		statuses, statusErr = service.GetStatusService().GetStatusesByProjectId(c.UserContext(), projectID)
	}()

	// Team members parallel fetch
	go func() {
		defer wg.Done()
		_, teamMembers, teamErr = service.GetProjectService().GetProjectWithTeam(c.UserContext(), projectID, user)
	}()

	wg.Wait()

	if statusErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": constant.ErrInternalServer})
	}
	if teamErr != nil {
		teamMembers = []models.User{}
	}

	var statusesWithIssues []response.StatusWithIssues
	for _, status := range statuses {
		if _, err := service.GetIssueService().GetIssuesByStatusID(c.UserContext(), status.ID); err != nil {
			fmt.Println(err)
		}

		statusesWithIssues = append(statusesWithIssues, response.StatusWithIssues{
			ID:        status.ID,
			ProjectID: status.ProjectID,
			Name:      status.Name,
			CreatedAt: status.CreatedAt,
			UpdatedAt: status.UpdatedAt,
			IssuesID:  status.IssueIDs,
		})
	}

	data := fiber.Map{
		"project":  project,
		"statutes": statusesWithIssues,
		"members":  teamMembers,
	}

	elapsed := time.Since(start)

	fmt.Println(elapsed)

	return c.JSON(fiber.Map{
		"message": constant.SuccessFetched,
		"data":    data,
	})
}

// @Summary Delete a member from a project by member ID
// @Description Deletes a member from a project using the member's ID.
// @Tags Projects
// @Produce json
// @Param memberId path string true "Member ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /projects/member/{memberId} [delete]
func DeleteMemberFromProjectByIdHandler(c *fiber.Ctx) error {

	memberId := c.Params("memberId")
	user, ok := utils.GetUserLocal(c)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	memberIdObj, err := primitive.ObjectIDFromHex(memberId)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": constant.ErrBadRequest})
	}

	err = service.GetProjectService().DeleteMemberFromProjectById(c.UserContext(), user.ID, memberIdObj)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": constant.ErrInternalServer,
		})
	}

	fmt.Print(err)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": constant.SuccessDeleted,
	})
}
