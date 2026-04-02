package service

import (
	"context"
	"fmt"
	"managify/database"
	"managify/internal/repository"

	"managify/models"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleService struct {
	roleRepo  repository.RoleRepository
	createLog func(*models.ProjectLog) error
}

var roleService *RoleService

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func GetRoleService() *RoleService {
	if roleService == nil {
		roleService = &RoleService{
			roleRepo:  repository.NewRoleRepository(database.DB),
			createLog: GetLogService().CreateLog,
		}
	}
	return roleService
}

func (s *RoleService) AddRole(userId primitive.ObjectID, projectId primitive.ObjectID, roleName string) (*models.Role, error) {

	ps := GetProjectService()

	projectValid, err := ps.IsProjectValid(projectId)
	if err != nil || !projectValid {
		return nil, err
	}

	role := &models.Role{
		ID:        primitive.NewObjectID(),
		UserID:    userId,
		ProjectID: projectId,
		RoleName:  roleName,
	}

	exists, err := ps.IsUserInProject(role.UserID, role.ProjectID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user is not part of the project")
	}

	roleRepo := s.roleRepo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = roleRepo.InsertOne(ctx, role)
	if err != nil {
		log.WithError(err).Error("failed to insert role")
		return nil, err
	}

	projectLogId := primitive.NewObjectID()
	projectLog := models.ProjectLog{
		ID:        projectLogId,
		ProjectID: role.ProjectID.Hex(),
		UserID:    userId.Hex(),
		Message:   "Role Has Been Assigned -> " + roleName,
		Timestamp: time.Now(),
	}
	if err := s.createLog(&projectLog); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *RoleService) DeleteRole(deleteId primitive.ObjectID) error {
	roleRepo := s.roleRepo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deletedCount, err := roleRepo.DeleteByID(ctx, deleteId)
	if err != nil {
		log.WithError(err).Error("failed to delete role")
		return err
	}

	if deletedCount == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

func (s *ProjectService) IsOwner(ownerId, projectId primitive.ObjectID) (bool, error) {
	projectRepo := s.projectRepo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	isOwner, err := projectRepo.VerifyProjectOwner(ctx, projectId, ownerId)
	if err != nil {
		return false, err
	}

	return isOwner, nil
}
