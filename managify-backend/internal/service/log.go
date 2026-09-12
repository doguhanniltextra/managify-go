package service

import (
	"context"
	"sync"
	"managify/database"
	"managify/internal/repository"
	"managify/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogService struct {
	logRepo repository.LogRepository
}

var logService *LogService
var logOnce sync.Once

func GetLogService() *LogService {
	logOnce.Do(func() {
		logService = &LogService{
			logRepo: repository.NewLogRepository(database.DB),
		}
	})
	return logService
}

func (s *LogService) CreateLog(ctx context.Context, projectLog *models.ProjectLog) error {
	projectLog.ID = primitive.NewObjectID()

	err := s.logRepo.InsertOne(ctx, projectLog)
	if err != nil {
		log.Errorf("Failed to insert log")
		return err
	}

	return nil
}

func (s *LogService) GetLogsByProjectID(ctx context.Context, projectID string) ([]models.ProjectLog, error) {
	logs, err := s.logRepo.FindByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *LogService) GetLogsByUserId(ctx context.Context, userID string) ([]models.ProjectLog, error) {
	logs, err := s.logRepo.GetRecentUserLogs(ctx, userID, 5)
	if err != nil {
		return nil, err
	}

	return logs, nil
}
