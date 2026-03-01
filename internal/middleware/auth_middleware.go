package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/mo/user-go-service/internal/service"
)

type contextKey string

const UserContextKey contextKey = "user"

// JWTMiddleware validates JWT tokens and adds user info to context
type JWTMiddleware struct {
	authService *service.AuthService
}

// NewJWTMiddleware creates a new instance of JWTMiddleware
func NewJWTMiddleware(authService *service.AuthService) *JWTMiddleware {
	return &JWTMiddleware{
		authService: authService,
	}
}

// Authenticate is a middleware that validates JWT tokens
func (m *JWTMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.authService.ValidateToken(tokenString)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole is a middleware that restricts access based on user role
func (m *JWTMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context (set by Authenticate middleware)
			claims, ok := r.Context().Value(UserContextKey).(*service.JWTClaims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				respondWithError(w, http.StatusForbidden, "Access denied: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// respondWithError sends an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
