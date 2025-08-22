package api

import (
	"bytes"
	"cli/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// APIService represents an API service with configuration
type APIService struct {
	config *config.Config
}

// JSONBinResponse represents the response from JSONBin API
type JSONBinResponse struct {
	Record   json.RawMessage `json:"record"`
	Metadata struct {
		ID        string `json:"id"`
		CreatedAt string `json:"createdAt"`
		Private   bool   `json:"private"`
	} `json:"metadata"`
}

// BinRecord represents a record to save locally
type BinRecord struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

// CreateBin sends a POST request to JSONBin API to create a new bin
func (api *APIService) CreateBin(fileContent string) (*JSONBinResponse, error) {
	url := "https://api.jsonbin.io/v3/b"

	// Get API key from config
	apiKey := api.GetConfigValue("KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found in configuration")
	}

	// Create request body
	body := bytes.NewBufferString(fileContent)

	// Create HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var jsonResponse JSONBinResponse
	if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &jsonResponse, nil
}

// GetBin sends a GET request to JSONBin API to retrieve a bin by ID
func (api *APIService) GetBin(binID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID)

	// Get API key from config
	apiKey := api.GetConfigValue("KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found in configuration")
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Master-Key", apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

// UpdateBin sends a PUT request to JSONBin API to update an existing bin
func (api *APIService) UpdateBin(binID, fileContent string) (*JSONBinResponse, error) {
	url := fmt.Sprintf("https://api.jsonbin.io/v3/b/%s", binID)

	// Get API key from config
	apiKey := api.GetConfigValue("KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found in configuration")
	}

	// Create request body
	body := bytes.NewBufferString(fileContent)

	// Create HTTP request
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var jsonResponse JSONBinResponse
	if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &jsonResponse, nil
}
