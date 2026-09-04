package service

import (
	"context"
	"fmt"
	"sync"
	"managify/database"
	"managify/internal/domain"
	"managify/internal/repository"

	"managify/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IssueService struct {
	issueRepo       repository.IssueRepository
	createLog       func(context.Context, *models.ProjectLog) error
	isProjectValid  func(context.Context, primitive.ObjectID) (bool, error)
	isUserInProject func(context.Context, primitive.ObjectID, primitive.ObjectID) (bool, error)
}

var issueService *IssueService
var issueOnce sync.Once

func GetIssueService() *IssueService {
	issueOnce.Do(func() {
		issueService = &IssueService{
			issueRepo:       repository.NewIssueRepository(database.DB),
			createLog:       GetLogService().CreateLog,
			isProjectValid:  GetProjectService().IsProjectValid,
			isUserInProject: GetProjectService().IsUserInProject,
		}
	})
	return issueService
}

func (s *IssueService) CreateIssue(ctx context.Context, issue *models.Issue, userID primitive.ObjectID) (*models.Issue, error) {

	// Project validation
	isProjectValid, err := s.isProjectValid(ctx, issue.ProjectID)
	if err != nil {
		return nil, err
	}
	if !isProjectValid {
		return nil, fmt.Errorf("%w: project is not valid", domain.ErrNotFound)
	}

	// User validation
	isUserInProject, err := s.isUserInProject(ctx, userID, issue.ProjectID)
	if err != nil {
		return nil, err
	}
	if !isUserInProject {
		return nil, fmt.Errorf("%w: user is not in project", domain.ErrForbidden)
	}

	issue.ID = primitive.NewObjectID()

	issueRepo := s.issueRepo

	if err := issueRepo.InsertOne(ctx, issue); err != nil {
		log.Errorf("Failed to insert issue into DB: %v", err)
		return nil, err
	}

	projectLogId := primitive.NewObjectID()
	projectLog := models.ProjectLog{
		ID:        projectLogId,
		ProjectID: issue.ProjectID.Hex(),
		UserID:    userID.Hex(),
		Message:   "Issue Has Been Created -> " + issue.Title,
		Timestamp: time.Now(),
	}
	if err := s.createLog(ctx, &projectLog); err != nil {
		return nil, err
	}

	return issue, nil
}
func (s *IssueService) DeleteIssue(ctx context.Context, issueID, userID primitive.ObjectID) error {
	issueRepo := s.issueRepo

	issue, err := issueRepo.FindByID(ctx, issueID)
	if err != nil {
		return err
	}
	if issue == nil {
		return fmt.Errorf("%w: issue not found", domain.ErrNotFound)
	}

	isUserInProject, err := s.isUserInProject(ctx, userID, issue.ProjectID)
	if err != nil {
		return err
	}
	if !isUserInProject {
		return fmt.Errorf("%w: user is not allowed to delete this issue", domain.ErrForbidden)
	}

	_, err = issueRepo.DeleteByID(ctx, issueID)
	if err != nil {
		log.Errorf("Failed to delete issue from DB: %v", err)
		return err
	}
	return nil
}
func (s *IssueService) GetIssuesByStatusID(ctx context.Context, statusID primitive.ObjectID) ([]*models.Issue, error) {

	issueRepo := s.issueRepo

	issues, err := issueRepo.FindByStatusID(ctx, statusID)
	if err != nil {
		return nil, err
	}
	return issues, nil
}
func (s *IssueService) UpdateIssueStatus(ctx context.Context, issueID, newStatusID, userID primitive.ObjectID) (*models.Issue, error) {
	issueRepo := s.issueRepo

	issue, err := issueRepo.FindByID(ctx, issueID)
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, fmt.Errorf("issue not found")
	}

	modifiedCount, err := issueRepo.UpdateStatus(ctx, issueID, newStatusID)
	if err != nil {
		return nil, fmt.Errorf("failed to update issue status: %w", err)
	}
	if modifiedCount == 0 {
		return nil, fmt.Errorf("no matching issue found to update")
	}

	projectLog := models.ProjectLog{
		ID:        primitive.NewObjectID(),
		ProjectID: issue.ProjectID.Hex(),
		UserID:    userID.Hex(),
		Message:   fmt.Sprintf("Issue '%s' status changed to new status", issue.Title),
		Timestamp: time.Now(),
	}
	if err := s.createLog(ctx, &projectLog); err != nil {
		return nil, err
	}

	issue.StatusID = newStatusID

	return issue, nil
}

func (s *IssueService) GetOncomingIssues(ctx context.Context, projectID primitive.ObjectID) ([]*models.Issue, error) {

	issueRepo := s.issueRepo
	currentTime := time.Now()
	threeDaysLater := currentTime.Add(72 * time.Hour)

	issues, err := issueRepo.FindOncomingIssues(ctx, projectID, currentTime, threeDaysLater)
	if err != nil {
		return nil, err
	}

	return issues, nil
}
