package service

import (
	"context"
	"managify/dto/request"
	"managify/internal/repository"
	"managify/models"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockProjectInviteRepository struct {
	repository.ProjectInviteRepository
	mock.Mock
}

func (m *mockProjectInviteRepository) CountByFilter(ctx context.Context, receiverID interface{}, projectID interface{}, statuses []string) (int64, error) {
	args := m.Called(ctx, receiverID, projectID, statuses)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockProjectInviteRepository) UpsertInvite(ctx context.Context, receiverID, projectID, senderID interface{}) (*models.ProjectInvite, error) {
	args := m.Called(ctx, receiverID, projectID, senderID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ProjectInvite), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProjectInviteRepository) FindInvitesByReceiverID(ctx context.Context, receiverID interface{}) ([]*models.ProjectInvite, error) {
	args := m.Called(ctx, receiverID)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.ProjectInvite), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProjectInviteRepository) UpdateStatus(ctx context.Context, inviteID, receiverID interface{}, newStatus string) (*models.ProjectInvite, error) {
	args := m.Called(ctx, inviteID, receiverID, newStatus)
	if args.Get(0) != nil {
		return args.Get(0).(*models.ProjectInvite), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockProjectRepository) AddUserToProject(ctx context.Context, projectID interface{}, userID interface{}) error {
	args := m.Called(ctx, projectID, userID)
	return args.Error(0)
}

func TestInviteService_CreateProjectInvite(t *testing.T) {
	mockInviteRepo := new(mockProjectInviteRepository)
	mockProjRepo := new(mockProjectRepository)
	mockUserRepo := new(mockUserRepository)

	senderId := primitive.NewObjectID()
	receiverId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	req := request.ProjectInviteRequest{
		Email:     "receiver@example.com",
		ProjectID: projectId.Hex(),
	}

	receiver := &models.User{
		ID:    receiverId,
		Email: req.Email,
	}

	project := &models.Project{
		ID:      projectId,
		TeamIDs: []primitive.ObjectID{},
	}

	invite := &models.ProjectInvite{
		ID: primitive.NewObjectID(),
	}

	mockUserRepo.On("FindByEmail", mock.Anything, req.Email).Return(receiver, nil)
	mockProjRepo.On("FindOneWithAccess", mock.Anything, projectId, senderId).Return(project, nil)
	mockInviteRepo.On("CountByFilter", mock.Anything, receiverId, projectId, mock.AnythingOfType("[]string")).Return(int64(0), nil)
	mockInviteRepo.On("UpsertInvite", mock.Anything, receiverId, projectId, senderId).Return(invite, nil)

	svc := &InviteService{
		inviteRepo:  mockInviteRepo,
		userRepo:    mockUserRepo,
		projectRepo: mockProjRepo,
		createLog:   func(ctx context.Context, pl *models.ProjectLog) error { return nil },
	}

	createdInvite, err := svc.CreateProjectInvite(context.Background(), senderId, req)

	assert.NoError(t, err)
	assert.NotNil(t, createdInvite)

	mockUserRepo.AssertExpectations(t)
	mockProjRepo.AssertExpectations(t)
	mockInviteRepo.AssertExpectations(t)
}

func TestInviteService_RespondProjectInvite(t *testing.T) {
	mockInviteRepo := new(mockProjectInviteRepository)
	mockProjRepo := new(mockProjectRepository)

	userId := primitive.NewObjectID()
	inviteId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	invite := &models.ProjectInvite{
		ID:        inviteId,
		ProjectID: projectId,
		Status:    "accepted",
	}

	mockInviteRepo.On("UpdateStatus", mock.Anything, inviteId, userId, "accepted").Return(invite, nil)
	mockProjRepo.On("AddUserToProject", mock.Anything, projectId, userId).Return(nil)

	svc := &InviteService{
		inviteRepo:  mockInviteRepo,
		projectRepo: mockProjRepo,
		createLog:   func(ctx context.Context, pl *models.ProjectLog) error { return nil },
	}

	res, err := svc.RespondProjectInvite(context.Background(), userId, inviteId, true)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "accepted", res.Status)

	mockInviteRepo.AssertExpectations(t)
	mockProjRepo.AssertExpectations(t)
}
