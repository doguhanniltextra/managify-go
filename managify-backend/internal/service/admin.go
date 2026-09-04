package service

import (
	"context"
	"fmt"
	"managify/internal/domain"
	"managify/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.FindAllUsers(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to find users")
		return nil, err
	}

	log.Infof("GetAllUsers succeeded, retrieved %d users", len(users))
	return users, nil
}

func (s *UserService) GetUserById(ctx context.Context, id string) (*models.User, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.WithError(err).Warnf("Invalid ObjectID format: %s", id)
		return nil, fmt.Errorf("%w: invalid id: %v", domain.ErrBadRequest, err)
	}

	user, err := s.userRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("%w: user not found", domain.ErrNotFound)
	}

	user.Password = "" // clear it implicitly as previously we used projection

	return user, nil
}

func (s *UserService) DeleteUserById(ctx context.Context, id string) (int64, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid ObjectID format: %v", domain.ErrBadRequest, err)
	}

	deletedCount, err := s.userRepo.DeleteByID(ctx, objID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %v", err)
	}

	if deletedCount == 0 {
		return 0, fmt.Errorf("%w: user not found", domain.ErrNotFound)
	}

	return deletedCount, nil
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	log.Debug("GetAllProjects called")

	projectRepo := s.projectRepo
	projects, err := projectRepo.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to find projects")
		return nil, err
	}

	return projects, nil
}

func (s *RoleService) GetAllRoles(ctx context.Context) ([]models.Role, error) {

	roleRepo := s.roleRepo
	roles, err := roleRepo.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to find roles")
		return nil, err
	}

	return roles, nil
}
