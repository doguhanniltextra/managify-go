package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"managify/database"
	"managify/dto/request"
	"managify/dto/response"
	"managify/internal/middleware"
	"managify/internal/repository"
	"net/smtp"
	"os"
	"time"

	"managify/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo        repository.UserRepository
	EncryptPassword func([]byte) ([]byte, error)
	CreateToken     func(*models.User) (string, error)
}

var userService *UserService

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			userRepo:        repository.NewUserRepository(database.DB),
			CreateToken:     middleware.CreateToken,
			EncryptPassword: encryptPassword,
		}
	}
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

func sendVerificationEmail(email, token string) error {

	from := os.Getenv("SMTP_FROM")
	pass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	fmt.Println(from, pass, smtpHost, smtpPort)

	to := email
	msg := "Subject: Email Verification\n" +
		"MIME-version: 1.0;\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\n\n" +
		"<html>" +
		"<body>" +
		"<h2>Verify Your Email</h2>" +
		"<p>Click the button below to verify your account:</p>" +
		"<a href='http://localhost:5173/verify?token=" + token + "' " +
		"style='display:inline-block;padding:10px 20px;background-color:#4CAF50;color:white;text-decoration:none;border-radius:5px;'>Verify Email</a>" +
		"<p>If you did not create an account, you can ignore this email.</p>" +
		"</body>" +
		"</html>"

	addr := smtpHost + ":" + smtpPort
	return smtp.SendMail(addr,
		smtp.PlainAuth("", from, pass, smtpHost),
		from, []string{to}, []byte(msg))
}

func (s *UserService) CreateUser(user *models.User) (*models.User, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hashedPassword, err := s.EncryptPassword([]byte(user.Password))
	if err != nil {
		log.Errorf("Password encryption failed: %v", err)
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
		log.Errorf("Failed to insert user into DB: %v", err)
		return nil, "", err
	}

	tokenString, err := s.CreateToken(user)
	if err != nil {
		log.Errorf("Failed to create JWT token: %v", err)
		return nil, "", err
	}

	go sendVerificationEmail(user.Email, user.VerificationToken)

	user.Password = ""

	return user, tokenString, nil
}

func (s *UserService) VerifyEmail(token string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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

func (s *UserService) Login(req *request.UserLoginRequest) (*response.UserLoginResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		log.Warnf("User not found: %s", req.Email)
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
		go sendVerificationEmail(user.Email, user.VerificationToken)
	}

	log.Infof("User logged in successfully: %s", req.Email)
	return resp, nil
}

func (s *UserService) IsUserValid(userId primitive.ObjectID) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByID(ctx, userId)
	if err != nil {
		log.WithError(err).Error("failed to fetch user")
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

func (s *UserService) GetUserByGivenId(givenId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(givenId)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, objID)
}
