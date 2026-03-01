package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mo/user-go-service/internal/models"
	"github.com/mo/user-go-service/internal/repository"
	"github.com/mo/user-go-service/internal/service"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	userRepo    repository.UserRepository
	authService *service.AuthService
}

// NewAuthHandler creates a new instance of AuthHandler
func NewAuthHandler(userRepo repository.UserRepository, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		authService: authService,
	}
}

// Register handles POST /auth/register
// @Summary Register a new user
// @Description Register a new user account with email, name, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Email == "" || req.Name == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email, name, and password are required")
		return
	}

	// Validate and set role (default to 'user' if not provided or invalid)
	role := req.Role
	if role == "" {
		role = models.RoleUser
	} else if role != models.RoleUser && role != models.RoleAdmin {
		respondWithError(w, http.StatusBadRequest, "Invalid role: must be 'user' or 'admin'")
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Create user
	user := &models.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Name:     req.Name,
		Role:     role,
		Password: hashedPassword,
	}

	if err := h.userRepo.Create(user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Email already exists")
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// Login handles POST /auth/login
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Compare passwords
	if err := h.authService.ComparePasswords(user.Password, req.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response := models.LoginResponse{
		Token: token,
		User:  user,
	}

	respondWithJSON(w, http.StatusOK, response)
}
