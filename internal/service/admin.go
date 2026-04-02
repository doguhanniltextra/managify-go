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

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetLevel(logrus.DebugLevel)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users, err := s.userRepo.FindAllUsers(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to find users")
		return nil, err
	}

	log.Infof("GetAllUsers succeeded, retrieved %d users", len(users))
	return users, nil
}

func (s *UserService) GetUserById(id string) (*models.User, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Warnf("Invalid ObjectID format: %s", id)
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := s.userRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	user.Password = "" // clear it implicitly as previously we used projection

	return user, nil
}

func (s *UserService) DeleteUserById(id string) (int64, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid ObjectID format: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deletedCount, err := s.userRepo.DeleteByID(ctx, objID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %v", err)
	}

	if deletedCount == 0 {
		return 0, fmt.Errorf("user not found")
	}

	return deletedCount, nil
}

func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	log.Debug("GetAllProjects called")

	projectRepo := repository.NewProjectRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	projects, err := projectRepo.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to find projects")
		return nil, err
	}

	return projects, nil
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {

	roleRepo := repository.NewRoleRepository(database.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	roles, err := roleRepo.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to find roles")
		return nil, err
	}

	return roles, nil
}
