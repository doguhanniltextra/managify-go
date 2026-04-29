package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"managify/dto/request"
	"managify/internal/notification"
	"managify/internal/repository"
	"managify/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockUserRepository embeds the original interface so we only mock what we need.
type mockUserRepository struct {
	repository.UserRepository
	mock.Mock
}

func (m *mockUserRepository) InsertOne(ctx context.Context, document *models.User) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func (m *mockUserRepository) FindByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) VerifyUser(ctx context.Context, userID interface{}) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id interface{}) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockNotificationProvider struct {
	notification.Provider
	mock.Mock
}

func (m *mockNotificationProvider) SendVerificationEmail(email, token string) error {
	args := m.Called(email, token)
	return args.Error(0)
}

func (m *mockNotificationProvider) SendInviteEmail(email, projectName, inviteToken string) error {
	args := m.Called(email, projectName, inviteToken)
	return args.Error(0)
}

func (m *mockNotificationProvider) SendPasswordResetEmail(email, token string) error {
	args := m.Called(email, token)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(mockUserRepository)
	mockNotifier := new(mockNotificationProvider)
	mockRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
	mockNotifier.On("SendVerificationEmail", mock.Anything, mock.Anything).Return(nil)

	svc := &UserService{
		userRepo: mockRepo,
		notifier: mockNotifier,
		EncryptPassword: func(pw []byte) ([]byte, error) {
			return []byte("hashed_" + string(pw)), nil
		},
		CreateToken: func(user *models.User) (string, error) {
			return "mock_jwt_token", nil
		},
	}

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	createdUser, token, err := svc.CreateUser(context.Background(), user)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "mock_jwt_token", token)
	assert.Equal(t, "", createdUser.Password) // password is removed before return
	assert.False(t, createdUser.IsVerified)   // should be initialized to false
	assert.NotEmpty(t, createdUser.ID)
	assert.NotEmpty(t, createdUser.VerificationToken)

	mockRepo.AssertExpectations(t)
}

func TestUserService_VerifyEmail(t *testing.T) {
	mockRepo := new(mockUserRepository)

	validToken := "valid_token"
	mockUser := &models.User{
		ID:                primitive.NewObjectID(),
		VerificationToken: validToken,
		IsVerified:        false,
	}

	mockRepo.On("FindByVerificationToken", mock.Anything, validToken).Return(mockUser, nil)
	mockRepo.On("VerifyUser", mock.Anything, mockUser.ID).Return(nil)

	svc := &UserService{
		userRepo: mockRepo,
	}

	verifiedUser, err := svc.VerifyEmail(context.Background(), validToken)
	assert.NoError(t, err)
	assert.NotNil(t, verifiedUser)
	assert.True(t, verifiedUser.IsVerified)
	assert.Empty(t, verifiedUser.VerificationToken)

	mockRepo.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(mockUserRepository)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:         primitive.NewObjectID(),
		Email:      "test@example.com",
		Password:   string(hashedPassword),
		IsVerified: true,
	}

	mockRepo.On("FindByEmail", mock.Anything, mockUser.Email).Return(mockUser, nil)

	svc := &UserService{
		userRepo: mockRepo,
		CreateToken: func(user *models.User) (string, error) {
			return "mock_token", nil
		},
	}

	req := &request.UserLoginRequest{
		Email:    mockUser.Email,
		Password: password,
	}

	resp, err := svc.Login(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, mockUser.Email, resp.Email)
	assert.Equal(t, "mock_token", resp.Token)

	mockRepo.AssertExpectations(t)
}

func TestUserService_IsUserValid(t *testing.T) {
	mockRepo := new(mockUserRepository)

	validId := primitive.NewObjectID()
	mockRepo.On("FindByID", mock.Anything, validId).Return(&models.User{}, nil)

	invalidId := primitive.NewObjectID()
	mockRepo.On("FindByID", mock.Anything, invalidId).Return(nil, nil)

	svc := &UserService{
		userRepo: mockRepo,
	}

	isValid, err := svc.IsUserValid(context.Background(), validId)
	assert.NoError(t, err)
	assert.True(t, isValid)

	isValid, err = svc.IsUserValid(context.Background(), invalidId)
	assert.NoError(t, err)
	assert.False(t, isValid)

	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByGivenId(t *testing.T) {
	mockRepo := new(mockUserRepository)

	validId := primitive.NewObjectID()
	mockUser := &models.User{ID: validId}

	mockRepo.On("FindByID", mock.Anything, validId).Return(mockUser, nil)

	svc := &UserService{
		userRepo: mockRepo,
	}

	user, err := svc.GetUserByGivenId(context.Background(), validId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, validId, user.ID)

	_, err = svc.GetUserByGivenId(context.Background(), "invalid-hex")
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}
