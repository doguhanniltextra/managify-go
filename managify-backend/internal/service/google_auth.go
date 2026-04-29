package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"managify/database"
	"managify/dto/response"
	"managify/internal/middleware"
	"managify/internal/repository"
	"managify/models"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleAuthService struct {
	userRepo repository.UserRepository
}

var googleAuthService *GoogleAuthService

func GetGoogleAuthService() *GoogleAuthService {
	if googleAuthService == nil {
		googleAuthService = &GoogleAuthService{
			userRepo: repository.NewUserRepository(database.DB),
		}
	}
	return googleAuthService
}

func (s *GoogleAuthService) GetGoogleAuthURL() string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("access_type", "offline") // refresh_token almak için
	params.Set("prompt", "consent")      // her seferinde izin ekranı göster (refresh_token garantisi)

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

func (s *GoogleAuthService) HandleCallback(ctx context.Context, code string) (*models.User, string, error) {

	tokenResp, err := s.exchangeCodeForToken(ctx, code)
	if err != nil {
		log.WithError(err).Error("Google token exchange failed")
		return nil, "", fmt.Errorf("google token exchange failed: %w", err)
	}

	userInfo, err := s.getUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		log.WithError(err).Error("Failed to get Google user info")
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}

	user, err := s.findOrCreateGoogleUser(ctx, userInfo)
	if err != nil {
		log.WithError(err).Error("Failed to find or create Google user")
		return nil, "", fmt.Errorf("failed to process user: %w", err)
	}

	token, err := middleware.CreateToken(user)
	if err != nil {
		log.WithError(err).Error("Failed to create JWT for Google user")
		return nil, "", fmt.Errorf("failed to create token: %w", err)
	}

	log.Infof("Google OAuth login successful: %s", user.Email)
	return user, token, nil
}

func (s *GoogleAuthService) exchangeCodeForToken(ctx context.Context, code string) (*response.GoogleTokenResponse, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Google token exchange returned %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("google returned status %d", resp.StatusCode)
	}

	var tokenResp response.GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// GET https://www.googleapis.com/oauth2/v3/userinfo
func (s *GoogleAuthService) getUserInfo(ctx context.Context, accessToken string) (*response.GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read userinfo response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Google userinfo returned %d: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("userinfo returned status %d", resp.StatusCode)
	}

	var userInfo response.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse userinfo response: %w", err)
	}

	return &userInfo, nil
}

func (s *GoogleAuthService) findOrCreateGoogleUser(ctx context.Context, info *response.GoogleUserInfo) (*models.User, error) {

	// 1. Repository üzerinden e-posta ile mevcut kullanıcıyı ara
	existingUser, err := s.userRepo.FindByEmail(ctx, info.Email)

	if err == nil && existingUser != nil {
		// Kullanıcı zaten var — google_id'yi güncelle (hesap bağlama)
		updateErr := s.userRepo.LinkGoogleID(ctx, existingUser.ID, info.Sub)
		if updateErr != nil {
			log.WithError(updateErr).Warn("Failed to link Google ID to existing user")
			// Link başarısız olsa bile login'e devam et
		}
		existingUser.GoogleID = info.Sub
		return existingUser, nil
	}

	// 2. Yoksa yeni kullanıcı oluştur
	newUser := &models.User{
		ID:           primitive.NewObjectID(),
		FullName:     info.Name,
		Email:        info.Email,
		Password:     "",
		GoogleID:     info.Sub,
		AuthProvider: "google",
		IsVerified:   true,
		IsAdmin:      false,
		ProjectSize:  0,
	}

	if insertErr := s.userRepo.InsertOne(ctx, newUser); insertErr != nil {
		return nil, fmt.Errorf("failed to create Google user: %w", insertErr)
	}

	subscription := models.Subscription{
		ID:                    primitive.NewObjectID(),
		UserID:                newUser.ID,
		PlanType:              models.PlanBasic,
		IsValid:               true,
		SubscriptionStartDate: time.Now(),
		SubscriptionEndDate:   time.Now(),
	}

	if _, err := GetSubscriptionService().CreateSubscription(ctx, &subscription); err != nil {
		log.WithError(err).Warn("Failed to create subscription for Google user")
	}

	log.Infof("New Google user created: %s (%s)", newUser.Email, newUser.ID.Hex())
	return newUser, nil
}
