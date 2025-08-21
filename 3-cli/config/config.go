package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Key string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Key: getEnv("KEY", "default_key"),
	}
}

// LoadFromEnvFile loads configuration from .env file and environment variables
// Note: This function requires the godotenv library. To install:
// go get github.com/joho/godotenv
func LoadFromEnvFile(envPath string) *Config {
	// For now, we'll just load from environment variables
	// To use .env file, uncomment the lines below and install godotenv

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return LoadConfig()
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetByKey returns a configuration value by key name
func (c *Config) GetByKey(key string) string {
	switch key {
	case "KEY":
		return c.Key
	default:
		log.Printf("Unknown configuration key: %s", key)
		return ""
	}
}

// Validate checks if all required configuration values are set
func (c *Config) Validate() error {
	// Add validation logic here if needed
	if c.Key == "" {
		log.Println("Warning: KEY is not set")
	}
	return nil
}
