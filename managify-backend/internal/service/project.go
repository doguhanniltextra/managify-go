package service

import (
	"context"
	"fmt"
	"managify/internal/repository"
	"managify/database"

	"managify/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectService struct {
	Collection string
}

var projectService *ProjectService

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetProjectService() *ProjectService {
	if projectService == nil {
		projectService = &ProjectService{Collection: "projects"}
	}
	return projectService
}

func (s *ProjectService) CreateProject(project *models.Project, user *models.User) (*models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userRepo := repository.NewUserRepository(database.DB)
	projectRepo := repository.NewProjectRepository(database.DB)
	subRepo := repository.NewSubscriptionRepository(database.DB)

	subscription, err := subRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check subscription: %w", err)
	}
	if subscription == nil || !subscription.IsValid {
		return nil, fmt.Errorf("no active subscription found")
	}

	if subscription.PlanType == models.PlanBasic && user.ProjectSize >= 3 {
		return nil, fmt.Errorf("plan limit reached: BASIC users can only create up to 3 projects")
	}

	err = userRepo.IncrementProjectSize(ctx, user.ID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to update project size: %w", err)
	}

	project.ID = primitive.NewObjectID()
	project.OwnerID = user.ID

	if err := projectRepo.InsertOne(ctx, project); err != nil {
		userRepo.IncrementProjectSize(ctx, user.ID, -1)
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}

	projectLog := models.ProjectLog{
		ID:        primitive.NewObjectID(),
		ProjectID: project.ID.Hex(),
		UserID:    user.ID.Hex(),
		Message:   "Project has been created",
		Timestamp: time.Now(),
	}
	if err := GetLogService().CreateLog(&projectLog); err != nil {
		return nil, err
	}

	return project, nil
}

func reduceProjectSize(ownerID primitive.ObjectID) error {
	log.Debugf("reduceProjectSize called for ownerID=%s", ownerID.Hex())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userRepo := repository.NewUserRepository(database.DB)
	err := userRepo.IncrementProjectSize(ctx, ownerID, -1)
	if err != nil {
		log.WithError(err).Error("Failed to update project_size")
		return fmt.Errorf("failed to update project size: %v", err)
	}

	return nil
}

func (s *ProjectService) DeleteProjectById(objID primitive.ObjectID, user *models.User) error {
	log.Debugf("DeleteProjectById called with objID=%s, userID=%s", objID.Hex(), user.ID.Hex())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRepo := repository.NewProjectRepository(database.DB)

	project, err := projectRepo.FindOneWithAccess(ctx, objID, user.ID)
	if err != nil {
		log.WithError(err).Error("Error finding project")
		return err
	}
	if project == nil {
		projectValid, _ := projectRepo.VerifyProject(ctx, objID)
		if projectValid && !user.IsAdmin {
			log.Warnf("Unauthorized delete attempt by user %s on project %s", user.ID.Hex(), objID.Hex())
			return fmt.Errorf("unauthorized: only owner or admin can delete")
		} else if !projectValid {
			log.Warnf("Project not found: %s", objID.Hex())
			return fmt.Errorf("project not found")
		}
	}

	deletedCount, err := projectRepo.DeleteByID(ctx, objID)
	if err != nil {
		log.WithError(err).Error("Failed to delete project")
		return err
	}

	if project != nil {
		reduceProjectSize(project.OwnerID)
	}

	log.Infof("Project deleted successfully: %s, deletedCount=%d", objID.Hex(), deletedCount)
	return nil
}

func (s *ProjectService) GetProject(projectID primitive.ObjectID, user *models.User) (*models.Project, error) {

	projectRepo := repository.NewProjectRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	project, err := projectRepo.FindOneWithAccess(ctx, projectID, user.ID)

	if err != nil {
		log.WithError(err).Error("failed to fetch project")
		return nil, err
	}
	if project == nil {
		log.Warnf("project not found or access denied for user %s", user.ID.Hex())
		return nil, fmt.Errorf("project not found or access denied")
	}

	return project, nil
}

func (s *ProjectService) IsProjectValid(projectID primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRepo := repository.NewProjectRepository(database.DB)
	valid, err := projectRepo.VerifyProject(ctx, projectID)
	if err != nil {
		log.WithError(err).Error("failed to fetch project")
		return false, err
	}

	return valid, nil
}

func (s *ProjectService) IsUserInProject(userID, projectID primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRepo := repository.NewProjectRepository(database.DB)
	inProject, err := projectRepo.CheckUserInProject(ctx, projectID, userID)
	if err != nil {
		log.WithError(err).Error("failed to check if user is in project")
		return false, err
	}

	return inProject, nil
}

func (s *ProjectService) GetProjectsByUserId(userIDHex string) ([]*models.Project, error) {
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectRepo := repository.NewProjectRepository(database.DB)
	projects, err := projectRepo.FindAllByUserID(ctx, userObjID)
	if err != nil {
		return nil, fmt.Errorf("decode projects failed: %w", err)
	}

	return projects, nil
}

func (s *ProjectService) GetProjectWithTeam(projectID primitive.ObjectID, user *models.User) (*models.Project, []models.User, error) {
	projectRepo := repository.NewProjectRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	project, err := projectRepo.FindOneWithAccess(ctx, projectID, user.ID)
	if err != nil {
		return nil, nil, err
	}
	if project == nil {
		return nil, nil, fmt.Errorf("project not found or access denied")
	}

	var teamMembers []models.User
	if len(project.TeamIDs) > 0 {
		userRepo := repository.NewUserRepository(database.DB)
		members, err := userRepo.FindUsersByIDs(ctx, project.TeamIDs)
		if err == nil {
			teamMembers = members
		}
	}

	return project, teamMembers, nil
}

func (s *ProjectService) DeleteMemberFromProjectById(userId, memberId primitive.ObjectID) error {
	projectRepo := repository.NewProjectRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	modifiedCount, err := projectRepo.RemoveMemberFromProject(ctx, userId, memberId)
	if err != nil {
		log.WithError(err).Error("failed to remove team member")
		return err
	}

	if modifiedCount == 0 {
		return fmt.Errorf("no modifications made, project not found or member not present") // could be either
	}
	return nil
}
