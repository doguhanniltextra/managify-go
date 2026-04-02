package service

import (
	"context"
	"managify/database"
	"managify/internal/repository"
	"managify/models"
	"time"

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

func (s *LogService) CreateLog(projectLog *models.ProjectLog) error {
	logRepo := repository.NewLogRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projectLog.ID = primitive.NewObjectID()

	err := logRepo.InsertOne(ctx, projectLog)
	if err != nil {
		log.Errorf("Failed to insert log")
		return err
	}

	return nil
}

func (s *LogService) GetLogsByProjectID(projectID string) ([]models.ProjectLog, error) {
	logRepo := repository.NewLogRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logs, err := logRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *LogService) GetLogsByUserId(userID string) ([]models.ProjectLog, error) {
	logRepo := repository.NewLogRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logs, err := logRepo.GetRecentUserLogs(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
