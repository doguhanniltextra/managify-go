package service

import (
	"context"
	"managify/models"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *mockSubscriptionRepository) InsertOne(ctx context.Context, document *models.Subscription) error {
	args := m.Called(ctx, document)
	return args.Error(0)
}

func TestSubscriptionService_GetByUserId(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	
	svc := &SubscriptionService{
		subscriptionRepo: mockRepo,
	}

	userId := primitive.NewObjectID()
	sub := &models.Subscription{
		ID: primitive.NewObjectID(),
	}

	mockRepo.On("FindByUserID", mock.Anything, userId).Return(sub, nil)

	res, err := svc.GetByUserId(userId.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, sub.ID, res.ID)

	mockRepo.AssertExpectations(t)
}

func TestSubscriptionService_CreateSubscription(t *testing.T) {
	mockRepo := new(mockSubscriptionRepository)
	
	svc := &SubscriptionService{
		subscriptionRepo: mockRepo,
	}

	sub := &models.Subscription{
		IsValid: true,
	}

	mockRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Subscription")).Return(nil)

	res, err := svc.CreateSubscription(sub)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.False(t, res.ID.IsZero())

	mockRepo.AssertExpectations(t)
}
