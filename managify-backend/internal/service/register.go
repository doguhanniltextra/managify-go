package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"managify/database"
	"managify/dto/request"
	"managify/dto/response"

	"managify/internal/config"
	"managify/internal/middleware"
	"managify/internal/notification"
	"managify/internal/repository"

	"managify/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo        repository.UserRepository
	notifier        notification.Provider
	EncryptPassword func([]byte) ([]byte, error)
	CreateToken     func(*models.User) (string, error)
}

var userService *UserService
var userOnce sync.Once

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetUserService() *UserService {
	userOnce.Do(func() {
		if userService == nil {
			cfg := config.LoadConfig()
			userService = &UserService{
				userRepo:        repository.NewUserRepository(database.DB),
				notifier:        notification.NewSMTPProvider(cfg),
				CreateToken:     middleware.CreateToken,
				EncryptPassword: encryptPassword,
			}
		}
	})
	return userService
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}


func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, string, error) {

	hashedPassword, err := s.EncryptPassword([]byte(user.Password))
	if err != nil {
		logrus.Errorf("Password encryption failed: %v", err)
		return nil, "", err
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	verifyToken, err := generateToken(32)
	if err != nil {
		return nil, "", err
	}
	user.VerificationToken = verifyToken
	user.IsVerified = false

	err = s.userRepo.InsertOne(ctx, user)
	if err != nil {
		logrus.Errorf("Failed to insert user into DB: %v", err)
		return nil, "", err
	}

	tokenString, err := s.CreateToken(user)
	if err != nil {
		logrus.Errorf("Failed to create JWT token: %v", err)
		return nil, "", err
	}

	go s.notifier.SendVerificationEmail(user.Email, user.VerificationToken)

	user.Password = ""

	return user, tokenString, nil
}

func (s *UserService) VerifyEmail(ctx context.Context, token string) (*models.User, error) {

	user, err := s.userRepo.FindByVerificationToken(ctx, token)
	if err != nil {
		fmt.Println("Error finding user with token:", err)
		return nil, err
	}

	err = s.userRepo.VerifyUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	user.IsVerified = true
	user.VerificationToken = ""

	return user, nil
}

func (s *UserService) Login(ctx context.Context, req *request.UserLoginRequest) (*response.UserLoginResponse, error) {

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		logrus.Warnf("User not found: %s", req.Email)
		return nil, fmt.Errorf("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	tokenString, err := s.CreateToken(user)
	if err != nil {
		return nil, fmt.Errorf("could not generate token")
	}

	resp := &response.UserLoginResponse{
		FullName: user.FullName,
		Email:    user.Email,
		Token:    tokenString,
	}
	if !user.IsVerified {
		go s.notifier.SendVerificationEmail(user.Email, user.VerificationToken)
	}

	logrus.Infof("User logged in successfully: %s", req.Email)
	return resp, nil
}

func (s *UserService) IsUserValid(ctx context.Context, userId primitive.ObjectID) (bool, error) {

	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch user")
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func encryptPassword(givenPassword []byte) (password []byte, error error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(givenPassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func (s *UserService) GetUserByGivenId(ctx context.Context, givenId string) (*models.User, error) {

	objID, err := primitive.ObjectIDFromHex(givenId)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, objID)
}
