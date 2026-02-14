package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port string
}

// Load returns the application configuration
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port: port,
	}
}
