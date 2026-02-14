package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mo/user-go-service/internal/config"
	"github.com/mo/user-go-service/internal/handlers"
	"github.com/mo/user-go-service/internal/repository"
	"github.com/mo/user-go-service/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize repository
	userRepo := repository.NewUserRepository()

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Setup router
	router := mux.NewRouter()
	
	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
