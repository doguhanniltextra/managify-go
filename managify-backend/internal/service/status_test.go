package service

import (
	"context"
	"managify/internal/repository"
	"managify/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockStatusRepository struct {
	repository.StatusRepository
	mock.Mock
}

func (m *mockStatusRepository) InsertOne(ctx context.Context, status *models.Status) error {
	args := m.Called(ctx, status)
	return args.Error(0)
}

func (m *mockStatusRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockStatusRepository) FindByProjectID(ctx context.Context, projectID interface{}) ([]*models.Status, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Status), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestStatusService_CreateStatus(t *testing.T) {
	mockStatusRepo := new(mockStatusRepository)

	svc := &StatusService{
		statusRepo:      mockStatusRepo,
		createLog:       func(ctx context.Context, pl *models.ProjectLog) error { return nil },
		isProjectValid:  func(ctx context.Context, pid primitive.ObjectID) (bool, error) { return true, nil },
		isUserInProject: func(ctx context.Context, uid, pid primitive.ObjectID) (bool, error) { return true, nil },
	}

	userId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	status := &models.Status{
		Name:      "Test Status",
		ProjectID: projectId,
		CreatorID: userId,
	}

	mockStatusRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Status")).Return(nil)

	created, err := svc.CreateStatus(context.Background(), status)
	assert.NoError(t, err)
	assert.NotNil(t, created)

	mockStatusRepo.AssertExpectations(t)
}

func TestStatusService_DeleteStatus(t *testing.T) {
	mockStatusRepo := new(mockStatusRepository)

	svc := &StatusService{
		statusRepo:      mockStatusRepo,
		isUserInProject: func(ctx context.Context, uid, pid primitive.ObjectID) (bool, error) { return true, nil },
	}

	statusId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()
	userId := primitive.NewObjectID()

	mockStatusRepo.On("DeleteByID", mock.Anything, statusId).Return(int64(1), nil)

	err := svc.DeleteStatus(context.Background(), statusId, projectId, userId)
	assert.NoError(t, err)

	mockStatusRepo.AssertExpectations(t)
}
