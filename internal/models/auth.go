package models

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"securePassword123"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  *User  `json:"user"`
}

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"securePassword123"`
	Role     string `json:"role,omitempty" example:"user"` // Optional, defaults to 'user'
}
