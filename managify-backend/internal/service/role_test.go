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

type mockRoleRepository struct {
	repository.RoleRepository
	mock.Mock
}

func (m *mockRoleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Role), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRoleRepository) InsertOne(ctx context.Context, role *models.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *mockRoleRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func TestRoleService_DeleteRole(t *testing.T) {
	mockRoleRepo := new(mockRoleRepository)
	roleId := primitive.NewObjectID()

	mockRoleRepo.On("DeleteByID", mock.Anything, roleId).Return(int64(1), nil)

	svc := &RoleService{
		roleRepo: mockRoleRepo,
	}

	err := svc.DeleteRole(roleId)
	assert.NoError(t, err)

	mockRoleRepo.AssertExpectations(t)
}

func TestRoleService_AddRole(t *testing.T) {
	mockRoleRepo := new(mockRoleRepository)
	mockProjectRepo := new(mockProjectRepository)

	userId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	// override global project service for the test
	originalProjectService := projectService
	defer func() { projectService = originalProjectService }()

	projectService = &ProjectService{
		projectRepo: mockProjectRepo,
	}

	mockProjectRepo.On("VerifyProject", mock.Anything, projectId).Return(true, nil)
	mockProjectRepo.On("CheckUserInProject", mock.Anything, projectId, userId).Return(true, nil)
	
	mockRoleRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Role")).Return(nil)

	svc := &RoleService{
		roleRepo: mockRoleRepo,
		createLog: func(pl *models.ProjectLog) error { return nil },
	}

	role, err := svc.AddRole(userId, projectId, "Admin")
	assert.NoError(t, err)
	assert.NotNil(t, role)
	assert.Equal(t, "Admin", role.RoleName)

	mockRoleRepo.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}

func TestProjectService_IsOwner(t *testing.T) {
	mockProjectRepo := new(mockProjectRepository)
	svc := &ProjectService{
		projectRepo: mockProjectRepo,
	}

	ownerId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	mockProjectRepo.On("VerifyProjectOwner", mock.Anything, projectId, ownerId).Return(true, nil)

	isOwner, err := svc.IsOwner(ownerId, projectId)
	assert.NoError(t, err)
	assert.True(t, isOwner)

	mockProjectRepo.AssertExpectations(t)
}
