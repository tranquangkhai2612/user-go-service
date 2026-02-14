package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mo/user-go-service/internal/models"
	"github.com/mo/user-go-service/internal/repository"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(req *models.CreateUserRequest) (*models.User, error)
	GetUser(id string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error)
	DeleteUser(id string) error
}

// userService is the implementation of UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Validate input
	if req.Email == "" || req.Name == "" {
		return nil, ErrInvalidInput
	}

	user := &models.User{
		ID:    uuid.New().String(),
		Email: req.Email,
		Name:  req.Name,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

// GetAllUsers retrieves all users
func (s *userService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAll()
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(id string, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
