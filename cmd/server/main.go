package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/mo/user-go-service/docs"
	"github.com/mo/user-go-service/internal/config"
	"github.com/mo/user-go-service/internal/handlers"
	"github.com/mo/user-go-service/internal/middleware"
	"github.com/mo/user-go-service/internal/repository"
	"github.com/mo/user-go-service/internal/service"
)

// @title User Service API
// @version 1.0
// @description A simple user management service with MySQL database and JWT authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to database")

	// Initialize repository with MySQL
	userRepo := repository.NewMySQLUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(cfg.JWTSecret, 24*time.Hour)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userRepo, authService)

	// Initialize middleware
	jwtMiddleware := middleware.NewJWTMiddleware(authService)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public auth routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected user routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(jwtMiddleware.Authenticate)

	// Admin-only routes (create, update, delete users)
	adminRoutes := protected.PathPrefix("").Subrouter()
	adminRoutes.Use(jwtMiddleware.RequireRole("admin"))
	adminRoutes.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// User routes (view users)
	protected.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	protected.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")

	// Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Swagger UI available at http://localhost:%s/swagger/", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
