package service

import (
	"context"
	"sync"
	"time"

	"managify/database"
	"managify/internal/repository"
	"managify/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository
}

var subscriptionService *SubscriptionService
var subscriptionOnce sync.Once

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetSubscriptionService() *SubscriptionService {
	subscriptionOnce.Do(func() {
		if subscriptionService == nil {
			subscriptionService = &SubscriptionService{
				subscriptionRepo: repository.NewSubscriptionRepository(database.DB),
			}
		}
	})
	return subscriptionService
}

func (s *SubscriptionService) GetByUserId(ctx context.Context, userIDHex string) (*models.Subscription, error) {
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, err
	}

	return s.subscriptionRepo.FindByUserID(ctx, userObjID)
}

func (s *SubscriptionService) CreateSubscription(ctx context.Context, subscription *models.Subscription) (*models.Subscription, error) {

	if subscription.ID.IsZero() {
		subscription.ID = primitive.NewObjectID()
	}

	err := s.subscriptionRepo.InsertOne(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) CreateDefaultSubscription(ctx context.Context, userID primitive.ObjectID) (*models.Subscription, error) {
	startDate := time.Now().UTC()
	endDate := startDate.AddDate(0, 1, 0) // Default to 1 month duration

	subscription := &models.Subscription{
		ID:                    primitive.NewObjectID(),
		UserID:                userID,
		SubscriptionStartDate: startDate,
		SubscriptionEndDate:   endDate,
		PlanType:              models.PlanBasic,
		IsValid:               true,
	}

	return s.CreateSubscription(ctx, subscription)
}
