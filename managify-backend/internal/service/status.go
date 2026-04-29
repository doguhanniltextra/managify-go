package service

import (
	"context"
	"fmt"
	"sync"
	"managify/database"
	"managify/internal/repository"
	"managify/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatusService struct {
	statusRepo      repository.StatusRepository
	createLog       func(context.Context, *models.ProjectLog) error
	isProjectValid  func(context.Context, primitive.ObjectID) (bool, error)
	isUserInProject func(context.Context, primitive.ObjectID, primitive.ObjectID) (bool, error)
}

var statusService *StatusService
var statusOnce sync.Once

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetStatusService() *StatusService {
	statusOnce.Do(func() {
		if statusService == nil {
			statusService = &StatusService{
				statusRepo:      repository.NewStatusRepository(database.DB),
				createLog:       GetLogService().CreateLog,
				isProjectValid:  GetProjectService().IsProjectValid,
				isUserInProject: GetProjectService().IsUserInProject,
			}
		}
	})
	return statusService
}

func (s *StatusService) CreateStatus(ctx context.Context, status *models.Status) (*models.Status, error) {
	projectValid, err := s.isProjectValid(ctx, status.ProjectID)
	if err != nil || !projectValid {
		return nil, err
	}

	exists, err := s.isUserInProject(ctx, status.CreatorID, status.ProjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user is not part of the project")
	}

	statusRepo := s.statusRepo

	status.CreatedAt = time.Now()
	status.UpdatedAt = time.Now()

	status.ID = primitive.NewObjectID()
	err = statusRepo.InsertOne(ctx, status)
	if err != nil {
		log.WithError(err).Error("failed to insert status")
		return nil, err
	}

	projectLogId := primitive.NewObjectID()
	projectLog := models.ProjectLog{
		ID:        projectLogId,
		ProjectID: status.ProjectID.Hex(),
		UserID:    status.ID.Hex(), // Is it intentional that user ID is status.ID? Left as is based on original
		Message:   "Status has been added -> " + status.Name,
		Timestamp: time.Now(),
	}
	if err := s.createLog(ctx, &projectLog); err != nil {
		return nil, err
	}
	return status, nil
}

func (s *StatusService) DeleteStatus(ctx context.Context, deleteId primitive.ObjectID, projectId primitive.ObjectID, userId primitive.ObjectID) error {
	statusRepo := s.statusRepo

	exists, err := s.isUserInProject(ctx, userId, projectId)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user is not part of the project")
	}

	deletedCount, err := statusRepo.DeleteByID(ctx, deleteId)
	if err != nil {
		log.WithError(err).Error("failed to delete status")
		return err
	}

	if deletedCount == 0 {
		return fmt.Errorf("status not found")
	}

	return nil
}

func (s *StatusService) GetStatusesByProjectId(ctx context.Context, projectID primitive.ObjectID) ([]*models.Status, error) {
	statusRepo := s.statusRepo

	statuses, err := statusRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
