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

type StatusService struct {
	Collection string
}

var statusService *StatusService

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetStatusService() *StatusService {
	if statusService == nil {
		statusService = &StatusService{Collection: "status"}
	}
	return statusService
}

func (s *StatusService) CreateStatus(status *models.Status) (*models.Status, error) {
	ps := GetProjectService()

	projectValid, err := ps.IsProjectValid(status.ProjectID)
	if err != nil || !projectValid {
		return nil, err
	}

	exists, err := ps.IsUserInProject(status.CreatorID, status.ProjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user is not part of the project")
	}

	statusRepo := repository.NewStatusRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
		UserID:    status.ID.Hex(),
		Message:   "Status has been added -> " + status.Name,
		Timestamp: time.Now(),
	}
	if err := GetLogService().CreateLog(&projectLog); err != nil {
		return nil, err
	}
	return status, nil
}

func (s *StatusService) DeleteStatus(deleteId primitive.ObjectID, projectId primitive.ObjectID, userId primitive.ObjectID) error {
	ps := GetProjectService()
	statusRepo := repository.NewStatusRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := ps.IsUserInProject(userId, projectId)
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

func (s *StatusService) GetStatusesByProjectId(projectID primitive.ObjectID) ([]*models.Status, error) {
	statusRepo := repository.NewStatusRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	statuses, err := statusRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}
