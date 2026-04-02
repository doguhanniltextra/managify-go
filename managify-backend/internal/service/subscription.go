package service

import (
	"context"
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

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetSubscriptionService() *SubscriptionService {
	if subscriptionService == nil {
		subscriptionService = &SubscriptionService{
			subscriptionRepo: repository.NewSubscriptionRepository(database.DB),
		}
	}
	return subscriptionService
}

func (s *SubscriptionService) GetByUserId(userIDHex string) (*models.Subscription, error) {
	userObjID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.subscriptionRepo.FindByUserID(ctx, userObjID)
}

func (s *SubscriptionService) CreateSubscription(subscription *models.Subscription) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if subscription.ID.IsZero() {
		subscription.ID = primitive.NewObjectID()
	}

	err := s.subscriptionRepo.InsertOne(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}
