package api

import (
	"cli/config"
	"fmt"
	"log"
)

// APIService represents an API service with configuration
type APIService struct {
	config *config.Config
}

// NewAPIService creates a new API service with the provided configuration
func NewAPIService(cfg *config.Config) *APIService {
	return &APIService{
		config: cfg,
	}
}

// GetConfig returns the current configuration
func (api *APIService) GetConfig() *config.Config {
	return api.config
}

// GetConfigValue returns a configuration value by key
func (api *APIService) GetConfigValue(key string) string {
	return api.config.GetByKey(key)
}

// Start initializes the API service with configuration
func (api *APIService) Start() error {
	// Validate configuration before starting
	if err := api.config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	log.Printf("API Service starting with configuration:")
	log.Printf("Key: %s", api.config.GetByKey("KEY"))

	// Add your API initialization logic here

	return nil
}

// Example API endpoint that uses configuration
func (api *APIService) HandleRequest() string {
	key := api.GetConfigValue("KEY")
	return fmt.Sprintf("Processing request with key: %s", key)
}
