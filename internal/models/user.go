package models

import "time"

// User represents a user in the system
type User struct {
	ID        string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	Password  string    `json:"-"` // Never expose password in JSON
	CreatedAt time.Time `json:"created_at" example:"2026-02-15T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2026-02-15T10:00:00Z"`
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"securePassword123"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Email string `json:"email,omitempty" example:"john.doe@example.com"`
	Name  string `json:"name,omitempty" example:"John Doe"`
}
