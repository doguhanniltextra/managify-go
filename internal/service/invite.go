package service

import (
	"context"
	"fmt"
	"managify/database"
	"managify/dto/request"
	"managify/internal/repository"
	"managify/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func CreateProjectInvite(senderID primitive.ObjectID, req request.ProjectInviteRequest) (*models.ProjectInvite, error) {
	inviteRepo := repository.NewProjectInviteRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	projectRepo := repository.NewProjectRepository(database.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	receiver, err := userRepo.FindByEmail(ctx, req.Email)
	if err != nil || receiver == nil {
		return nil, fmt.Errorf("receiver not found")
	}

	projectID, err := primitive.ObjectIDFromHex(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project ID")
	}

	project, err := projectRepo.FindOneWithAccess(ctx, projectID, senderID)
	if err != nil || project == nil {
		return nil, fmt.Errorf("project not found or access denied")
	}

	for _, member := range project.TeamIDs {
		if member == receiver.ID {
			return nil, fmt.Errorf("user is already a member of this project")
		}
	}

	count, err := inviteRepo.CountByFilter(ctx, receiver.ID, projectID, []string{"pending", "accepted"})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("invite already sent to this user")
	}

	invite, err := inviteRepo.UpsertInvite(ctx, receiver.ID, projectID, senderID)
	if err != nil {
		return nil, fmt.Errorf("invite already exists or could not be created")
	}

	projectLogId := primitive.NewObjectID()
	projectLog := models.ProjectLog{
		ID:        projectLogId,
		ProjectID: projectID.Hex(),
		UserID:    senderID.Hex(),
		Message:   "Invite has been sent to " + req.Email,
		Timestamp: time.Now(),
	}

	if err := GetLogService().CreateLog(&projectLog); err != nil {
		return nil, err
	}

	return invite, nil
}

type ProjectInviteFull struct {
	ID        primitive.ObjectID `json:"id"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"createdAt"`
	Project   models.Project     `json:"project"`
	Sender    models.User        `json:"sender"`
	Receiver  models.User        `json:"receiver"`
}

func GetProjectInvites(receiverID primitive.ObjectID) ([]*ProjectInviteFull, error) {
	inviteRepo := repository.NewProjectInviteRepository(database.DB)
	userRepo := repository.NewUserRepository(database.DB)
	projectRepo := repository.NewProjectRepository(database.DB)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	invites, err := inviteRepo.FindInvitesByReceiverID(ctx, receiverID)
	if err != nil {
		return nil, err
	}

	var result []*ProjectInviteFull

	for _, invite := range invites {
		project, err := projectRepo.FindOneWithAccess(ctx, invite.ProjectID, receiverID) // since user receives it
		if err != nil || project == nil {
			project = &models.Project{Name: "Project"} // fallback
		}

		sender, err := userRepo.FindByID(ctx, invite.SenderID)
		if err != nil || sender == nil {
			sender = &models.User{FullName: "Someone"} // fallback
		}

		receiver, err := userRepo.FindByID(ctx, invite.ReceiverID)
		if err != nil || receiver == nil {
			receiver = &models.User{FullName: "Unknown"}
		}

		result = append(result, &ProjectInviteFull{
			ID:        invite.ID,
			Status:    invite.Status,
			CreatedAt: invite.CreatedAt,
			Project:   *project,
			Sender:    *sender,
			Receiver:  *receiver,
		})
	}

	return result, nil
}

func RespondProjectInvite(userID, inviteID primitive.ObjectID, accept bool) (*models.ProjectInvite, error) {
	log.Debugf("RespondProjectInvite called with userID=%s, inviteID=%s, accept=%v", userID.Hex(), inviteID.Hex(), accept)

	inviteRepo := repository.NewProjectInviteRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := "declined"
	if accept {
		status = "accepted"
	}
	log.Debugf("Setting invite status to: %s", status)

	invite, err := inviteRepo.UpdateStatus(ctx, inviteID, userID, status)
	if err != nil {
		log.WithError(err).Warnf("Invite not found or already handled for inviteID=%s, userID=%s", inviteID.Hex(), userID.Hex())
		return nil, fmt.Errorf("invite not found or already handled")
	}
	log.Infof("Invite updated successfully: %+v", *invite)

	if accept {
		if err := addUserToProject(invite.ProjectID, userID); err != nil {
			log.WithError(err).Errorf("Failed to add user to project: projectID=%s, userID=%s", invite.ProjectID.Hex(), userID.Hex())
			return nil, fmt.Errorf("failed to add user to project: %v", err)
		}
		projectLogId := primitive.NewObjectID()
		projectLog := models.ProjectLog{
			ID:        projectLogId,
			ProjectID: invite.ProjectID.Hex(),
			UserID:    userID.Hex(),
			Message:   "Invite has been accepted",
			Timestamp: time.Now(),
		}
		if err := GetLogService().CreateLog(&projectLog); err != nil {
			return nil, err
		}
		log.Infof("User %s added to project %s team", userID.Hex(), invite.ProjectID.Hex())
	}

	return invite, nil
}

func addUserToProject(projectID, userID primitive.ObjectID) error {
	log.Debugf("addUserToProject called with projectID=%s, userID=%s", projectID.Hex(), userID.Hex())
	projectRepo := repository.NewProjectRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := projectRepo.AddUserToProject(ctx, projectID, userID)
	if err != nil {
		log.WithError(err).Error("Failed to update project team")
		return err
	}
	return nil
}
