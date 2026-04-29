package service

import (
	"context"
	"fmt"
	"managify/database"
	"managify/internal/repository"

	"managify/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProjectService struct {
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
	subRepo     repository.SubscriptionRepository
	createLog   func(context.Context, *models.ProjectLog) error
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
		projectService = &ProjectService{
			projectRepo: repository.NewProjectRepository(database.DB),
			userRepo:    repository.NewUserRepository(database.DB),
			subRepo:     repository.NewSubscriptionRepository(database.DB),
			createLog:   GetLogService().CreateLog,
		}
	}
	return projectService
}

func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project, user *models.User) (*models.Project, error) {

	userRepo := s.userRepo
	projectRepo := s.projectRepo
	subRepo := s.subRepo

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
	if err := s.createLog(ctx, &projectLog); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) reduceProjectSize(ctx context.Context, ownerID primitive.ObjectID) error {
	log.Debugf("reduceProjectSize called for ownerID=%s", ownerID.Hex())

	userRepo := s.userRepo
	err := userRepo.IncrementProjectSize(ctx, ownerID, -1)
	if err != nil {
		log.WithError(err).Error("Failed to update project_size")
		return fmt.Errorf("failed to update project size: %v", err)
	}

	return nil
}

func (s *ProjectService) DeleteProjectById(ctx context.Context, objID primitive.ObjectID, user *models.User) error {
	log.Debugf("DeleteProjectById called with objID=%s, userID=%s", objID.Hex(), user.ID.Hex())

	projectRepo := s.projectRepo

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
		s.reduceProjectSize(ctx, project.OwnerID)
	}

	log.Infof("Project deleted successfully: %s, deletedCount=%d", objID.Hex(), deletedCount)
	return nil
}

func (s *ProjectService) GetProject(ctx context.Context, projectID primitive.ObjectID, user *models.User) (*models.Project, error) {

	projectRepo := s.projectRepo

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

func (s *ProjectService) IsProjectValid(ctx context.Context, projectID primitive.ObjectID) (bool, error) {

	projectRepo := s.projectRepo
	valid, err := projectRepo.VerifyProject(ctx, projectID)
	if err != nil {
		log.WithError(err).Error("failed to fetch project")
		return false, err
	}

	return valid, nil
}

func (s *ProjectService) IsUserInProject(ctx context.Context, userID, projectID primitive.ObjectID) (bool, error) {

	projectRepo := s.projectRepo
	inProject, err := projectRepo.CheckUserInProject(ctx, projectID, userID)
	if err != nil {
		log.WithError(err).Error("failed to check if user is in project")
		return false, err
	}

	return inProject, nil
}

func (s *ProjectService) GetProjectsByUserId(ctx context.Context, userIDHex string) ([]*models.Project, error) {
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)

	projectRepo := s.projectRepo
	projects, err := projectRepo.FindAllByUserID(ctx, userObjID)
	if err != nil {
		return nil, fmt.Errorf("decode projects failed: %w", err)
	}

	return projects, nil
}

func (s *ProjectService) GetProjectWithTeam(ctx context.Context, projectID primitive.ObjectID, user *models.User) (*models.Project, []models.User, error) {
	projectRepo := s.projectRepo

	project, err := projectRepo.FindOneWithAccess(ctx, projectID, user.ID)
	if err != nil {
		return nil, nil, err
	}
	if project == nil {
		return nil, nil, fmt.Errorf("project not found or access denied")
	}

	var teamMembers []models.User
	if len(project.TeamIDs) > 0 {
		userRepo := s.userRepo
		members, err := userRepo.FindUsersByIDs(ctx, project.TeamIDs)
		if err == nil {
			teamMembers = members
		}
	}

	return project, teamMembers, nil
}

func (s *ProjectService) DeleteMemberFromProjectById(ctx context.Context, userId, memberId primitive.ObjectID) error {
	projectRepo := s.projectRepo

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
