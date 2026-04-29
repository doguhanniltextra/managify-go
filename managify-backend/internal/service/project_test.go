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

type mockProjectRepository struct {
	repository.ProjectRepository
	mock.Mock
}

func (m *mockProjectRepository) FindOneWithAccess(ctx context.Context, projectID interface{}, userID interface{}) (*models.Project, error) {
	args := m.Called(ctx, projectID, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Project), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProjectRepository) InsertOne(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *mockProjectRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockProjectRepository) VerifyProject(ctx context.Context, projectID interface{}) (bool, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockProjectRepository) CheckUserInProject(ctx context.Context, projectID interface{}, userID interface{}) (bool, error) {
	args := m.Called(ctx, projectID, userID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *mockProjectRepository) VerifyProjectOwner(ctx context.Context, projectID interface{}, ownerID interface{}) (bool, error) {
	args := m.Called(ctx, projectID, ownerID)
	return args.Get(0).(bool), args.Error(1)
}

type mockSubscriptionRepository struct {
	repository.SubscriptionRepository
	mock.Mock
}

func (m *mockSubscriptionRepository) FindByUserID(ctx context.Context, userID interface{}) (*models.Subscription, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Subscription), args.Error(1)
	}
	return nil, args.Error(1)
}

// Add method IncrementProjectSize to mockUserRepository
func (m *mockUserRepository) IncrementProjectSize(ctx context.Context, userID interface{}, amount int) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func TestProjectService_CreateProject(t *testing.T) {
	mockProjectRepo := new(mockProjectRepository)
	mockUserRepo := new(mockUserRepository)
	mockSubRepo := new(mockSubscriptionRepository)

	user := &models.User{
		ID: primitive.NewObjectID(),
		ProjectSize: 2,
	}

	sub := &models.Subscription{
		IsValid: true,
		PlanType: models.PlanBasic,
	}

	mockSubRepo.On("FindByUserID", mock.Anything, user.ID).Return(sub, nil)
	mockUserRepo.On("IncrementProjectSize", mock.Anything, user.ID, 1).Return(nil)
	mockProjectRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Project")).Return(nil)

	svc := &ProjectService{
		projectRepo: mockProjectRepo,
		userRepo:    mockUserRepo,
		subRepo:     mockSubRepo,
		createLog:   func(ctx context.Context, pl *models.ProjectLog) error { return nil },
	}

	project := &models.Project{
		Name: "Test Project",
	}

	createdProject, err := svc.CreateProject(context.Background(), project, user)
	
	assert.NoError(t, err)
	assert.NotNil(t, createdProject)
	assert.Equal(t, createdProject.OwnerID, user.ID)

	mockSubRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestProjectService_DeleteProjectById(t *testing.T) {
	mockProjectRepo := new(mockProjectRepository)
	mockUserRepo := new(mockUserRepository)

	user := &models.User{
		ID: primitive.NewObjectID(),
	}
	projectId := primitive.NewObjectID()

	existingProject := &models.Project{
		ID: projectId,
		OwnerID: user.ID,
	}

	mockProjectRepo.On("FindOneWithAccess", mock.Anything, projectId, user.ID).Return(existingProject, nil)
	mockProjectRepo.On("DeleteByID", mock.Anything, projectId).Return(int64(1), nil)
	mockUserRepo.On("IncrementProjectSize", mock.Anything, user.ID, -1).Return(nil)

	svc := &ProjectService{
		projectRepo: mockProjectRepo,
		userRepo:    mockUserRepo,
	}

	err := svc.DeleteProjectById(context.Background(), projectId, user)
	
	assert.NoError(t, err)

	mockProjectRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
