package service

import (
	"context"
	"managify/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *mockUserRepository) FindAllUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := new(mockUserRepository)

	users := []models.User{
		{ID: primitive.NewObjectID(), Email: "user1@example.com"},
		{ID: primitive.NewObjectID(), Email: "user2@example.com"},
	}

	mockRepo.On("FindAllUsers", mock.Anything).Return(users, nil)

	svc := &UserService{
		userRepo: mockRepo,
	}

	result, err := svc.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserById_Admin(t *testing.T) {
	mockRepo := new(mockUserRepository)

	validId := primitive.NewObjectID()
	mockUser := &models.User{
		ID:       validId,
		Password: "secretpassword",
	}

	mockRepo.On("FindByID", mock.Anything, validId).Return(mockUser, nil)

	svc := &UserService{
		userRepo: mockRepo,
	}

	resultUser, err := svc.GetUserById(validId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, validId, resultUser.ID)
	// Make sure password is explicitly cleared
	assert.Equal(t, "", resultUser.Password)

	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUserById(t *testing.T) {
	mockRepo := new(mockUserRepository)

	validId := primitive.NewObjectID()
	mockRepo.On("DeleteByID", mock.Anything, validId).Return(int64(1), nil)

	svc := &UserService{
		userRepo: mockRepo,
	}

	deletedCount, err := svc.DeleteUserById(validId.Hex())
	assert.NoError(t, err)
	assert.Equal(t, int64(1), deletedCount)

	mockRepo.AssertExpectations(t)
}
