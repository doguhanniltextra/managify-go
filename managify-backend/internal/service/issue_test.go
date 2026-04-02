package service

import (
	"context"
	"managify/internal/repository"
	"managify/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockIssueRepository struct {
	repository.IssueRepository
	mock.Mock
}

func (m *mockIssueRepository) InsertOne(ctx context.Context, issue *models.Issue) error {
	args := m.Called(ctx, issue)
	return args.Error(0)
}

func (m *mockIssueRepository) DeleteByID(ctx context.Context, id interface{}) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockIssueRepository) FindByID(ctx context.Context, id interface{}) (*models.Issue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Issue), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockIssueRepository) FindByStatusID(ctx context.Context, statusID interface{}) ([]*models.Issue, error) {
	args := m.Called(ctx, statusID)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Issue), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockIssueRepository) UpdateStatus(ctx context.Context, issueID interface{}, newStatusID interface{}) (int64, error) {
	args := m.Called(ctx, issueID, newStatusID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockIssueRepository) FindOncomingIssues(ctx context.Context, projectID interface{}, currentTime time.Time, limitTime time.Time) ([]*models.Issue, error) {
	args := m.Called(ctx, projectID, currentTime, limitTime)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Issue), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestIssueService_CreateIssue(t *testing.T) {
	mockIssueRepo := new(mockIssueRepository)

	svc := &IssueService{
		issueRepo:       mockIssueRepo,
		createLog:       func(pl *models.ProjectLog) error { return nil },
		isProjectValid:  func(pid primitive.ObjectID) (bool, error) { return true, nil },
		isUserInProject: func(uid, pid primitive.ObjectID) (bool, error) { return true, nil },
	}

	userId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	issue := &models.Issue{
		Title:     "Test Issue",
		ProjectID: projectId,
	}

	mockIssueRepo.On("InsertOne", mock.Anything, mock.AnythingOfType("*models.Issue")).Return(nil)

	created, err := svc.CreateIssue(issue, userId)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, issue.Title, created.Title)

	mockIssueRepo.AssertExpectations(t)
}

func TestIssueService_DeleteIssue(t *testing.T) {
	mockIssueRepo := new(mockIssueRepository)

	svc := &IssueService{
		issueRepo:       mockIssueRepo,
		isUserInProject: func(uid, pid primitive.ObjectID) (bool, error) { return true, nil },
	}

	userId := primitive.NewObjectID()
	issueId := primitive.NewObjectID()
	projectId := primitive.NewObjectID()

	existingIssue := &models.Issue{
		ID:        issueId,
		ProjectID: projectId,
	}

	mockIssueRepo.On("FindByID", mock.Anything, issueId).Return(existingIssue, nil)
	mockIssueRepo.On("DeleteByID", mock.Anything, issueId).Return(int64(1), nil)

	err := svc.DeleteIssue(issueId, userId)
	assert.NoError(t, err)

	mockIssueRepo.AssertExpectations(t)
}
