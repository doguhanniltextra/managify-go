package service

import (
	"context"
	"managify/database"
	"managify/internal/repository"
	"managify/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogService struct {
	Collection string
}

var logService *LogService

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetLogService() *LogService {
	if logService == nil {
		logService = &LogService{Collection: "logs"}
	}
	return logService
}

func (s *LogService) CreateLog(ctx context.Context, projectLog *models.ProjectLog) error {
	logRepo := repository.NewLogRepository(database.DB)

	projectLog.ID = primitive.NewObjectID()

	err := logRepo.InsertOne(ctx, projectLog)
	if err != nil {
		log.Errorf("Failed to insert log")
		return err
	}

	return nil
}

func (s *LogService) GetLogsByProjectID(ctx context.Context, projectID string) ([]models.ProjectLog, error) {
	logRepo := repository.NewLogRepository(database.DB)

	logs, err := logRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *LogService) GetLogsByUserId(ctx context.Context, userID string) ([]models.ProjectLog, error) {
	logRepo := repository.NewLogRepository(database.DB)

	logs, err := logRepo.GetRecentUserLogs(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
