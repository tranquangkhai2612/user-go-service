package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/mo/user-go-service/internal/models"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetAll() ([]*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*models.User
}

// NewUserRepository creates a new instance of InMemoryUserRepository
func NewUserRepository() UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*models.User),
	}
}

// Create adds a new user to the repository
func (r *InMemoryUserRepository) Create(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return ErrUserAlreadyExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetAll retrieves all users
func (r *InMemoryUserRepository) GetAll() ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// Delete removes a user from the repository
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}
