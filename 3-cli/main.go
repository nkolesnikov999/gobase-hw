package main

import (
	"cli/api"
	"cli/config"
	"cli/file"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// BinRecord represents a single bin record
type BinRecord struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BinsList represents a collection of bin records
type BinsList struct {
	Bins []BinRecord `json:"bins"`
}

func main() {
	// Define CLI flags
	var (
		create   = flag.Bool("create", false, "Create a new bin")
		list     = flag.Bool("list", false, "List all created bins")
		get      = flag.Bool("get", false, "Get a bin by ID")
		filename = flag.String("file", "", "Path to the JSON file to upload")
		binName  = flag.String("name", "", "Name for the bin")
		binID    = flag.String("id", "", "ID of the bin to retrieve")
	)

	flag.Parse()

	// Check command flags
	if !*create && !*list && !*get {
		fmt.Println("Usage:")
		fmt.Println("  ./cli --create --file=<path> --name=<name>  # Create a new bin")
		fmt.Println("  ./cli --list                                # List all created bins")
		fmt.Println("  ./cli --get --id=<bin_id>                   # Get bin data by ID")
		fmt.Println("Examples:")
		fmt.Println("  ./cli --create --file=bins_data.json --name=my-bin")
		fmt.Println("  ./cli --list")
		fmt.Println("  ./cli --get --id=68a83742ae596e708fd0f72c")
		os.Exit(1)
	}

	// Handle list command
	if *list {
		listBins()
		return
	}

	// Handle get command
	if *get {
		if *binID == "" {
			log.Fatal("Error: --id flag is required when using --get")
		}
		getBin(*binID)
		return
	}

	// Validate required flags
	if *filename == "" {
		log.Fatal("Error: --file flag is required")
	}

	if *binName == "" {
		log.Fatal("Error: --name flag is required")
	}

	// Load configuration
	cfg := config.LoadFromEnvFile(".env")

	// Check if API key is configured
	apiKey := cfg.GetByKey("KEY")
	if apiKey == "" {
		log.Fatal("Error: API key not found. Please create .env file with KEY=<your_api_key>")
	}

	fmt.Printf("Creating bin '%s' from file '%s'...\n", *binName, *filename)

	// Read file content
	fileContent, err := file.ReadFile(*filename)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Validate that file contains valid JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(fileContent), &jsonData); err != nil {
		log.Fatalf("Error: File does not contain valid JSON: %v", err)
	}

	// Create API service
	apiService := api.NewAPIService(cfg)

	// Create bin via JSONBin API
	response, err := apiService.CreateBin(fileContent)
	if err != nil {
		log.Fatalf("Error creating bin: %v", err)
	}

	fmt.Printf("Bin created successfully!\n")
	fmt.Printf("ID: %s\n", response.Metadata.ID)
	fmt.Printf("Created At: %s\n", response.Metadata.CreatedAt)
	fmt.Printf("Private: %t\n", response.Metadata.Private)

	// Save bin to local storage
	if err := saveBinToList(response.Metadata.ID, *binName); err != nil {
		log.Fatalf("Error saving bin to local storage: %v", err)
	}

	fmt.Printf("Bin information saved to bins.json\n")
}

const binsFileName = "bins.json"

// loadBinsList loads the existing bins list from file
func loadBinsList() (*BinsList, error) {
	// Check if file exists
	if _, err := os.Stat(binsFileName); os.IsNotExist(err) {
		// File doesn't exist, return empty list
		return &BinsList{Bins: []BinRecord{}}, nil
	}

	// Read file content
	content, err := file.ReadFile(binsFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read bins file: %w", err)
	}

	// Parse JSON
	var binsList BinsList
	if err := json.Unmarshal([]byte(content), &binsList); err != nil {
		return nil, fmt.Errorf("failed to parse bins file: %w", err)
	}

	return &binsList, nil
}

// saveBinToList adds a new bin to the list and saves it to file
func saveBinToList(id, name string) error {
	// Load existing list
	binsList, err := loadBinsList()
	if err != nil {
		return fmt.Errorf("failed to load existing bins: %w", err)
	}

	// Check if bin already exists
	for _, bin := range binsList.Bins {
		if bin.ID == id {
			return fmt.Errorf("bin with ID %s already exists", id)
		}
	}

	// Add new bin
	newBin := BinRecord{
		ID:   id,
		Name: name,
	}
	binsList.Bins = append(binsList.Bins, newBin)

	// Convert to JSON
	data, err := json.MarshalIndent(binsList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bins list: %w", err)
	}

	// Save to file
	if err := file.WriteFile(binsFileName, string(data)); err != nil {
		return fmt.Errorf("failed to save bins file: %w", err)
	}

	return nil
}

// listBins displays all created bins
func listBins() {
	binsList, err := loadBinsList()
	if err != nil {
		log.Fatalf("Error loading bins list: %v", err)
	}

	if len(binsList.Bins) == 0 {
		fmt.Println("No bins created yet.")
		fmt.Println("Use './cli --create --file=<path> --name=<name>' to create a bin.")
		return
	}

	fmt.Printf("Created bins (%d total):\n", len(binsList.Bins))
	fmt.Println("─────────────────────────────────────────────────────────────────")
	fmt.Printf("%-30s %s\n", "ID", "Name")
	fmt.Println("─────────────────────────────────────────────────────────────────")

	for _, bin := range binsList.Bins {
		fmt.Printf("%-30s %s\n", bin.ID, bin.Name)
	}
}

// getBin retrieves and displays a bin by ID
func getBin(binID string) {
	// Load configuration
	cfg := config.LoadFromEnvFile(".env")

	// Check if API key is configured
	apiKey := cfg.GetByKey("KEY")
	if apiKey == "" {
		log.Fatal("Error: API key not found. Please create .env file with KEY=<your_api_key>")
	}

	fmt.Printf("Getting bin with ID: %s...\n", binID)

	// Create API service
	apiService := api.NewAPIService(cfg)

	// Get bin via JSONBin API
	response, err := apiService.GetBin(binID)
	if err != nil {
		log.Fatalf("Error getting bin: %v", err)
	}

	// Extract and display the bin data
	fmt.Println("\nBin data retrieved successfully!")
	fmt.Println("═══════════════════════════════════════════════════════════════════")

	// Display metadata if available
	if metadata, ok := response["metadata"]; ok {
		if metadataMap, ok := metadata.(map[string]interface{}); ok {
			fmt.Println("Metadata:")
			if id, ok := metadataMap["id"]; ok {
				fmt.Printf("  ID: %v\n", id)
			}
			if createdAt, ok := metadataMap["createdAt"]; ok {
				fmt.Printf("  Created At: %v\n", createdAt)
			}
			if private, ok := metadataMap["private"]; ok {
				fmt.Printf("  Private: %v\n", private)
			}
			fmt.Println()
		}
	}

	// Display the actual record/data
	if record, ok := response["record"]; ok {
		fmt.Println("Record Data:")
		recordJSON, err := json.MarshalIndent(record, "", "  ")
		if err != nil {
			fmt.Printf("  %v\n", record)
		} else {
			fmt.Printf("%s\n", recordJSON)
		}
	} else {
		// If no record field, display the entire response
		fmt.Println("Response Data:")
		responseJSON, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			fmt.Printf("  %v\n", response)
		} else {
			fmt.Printf("%s\n", responseJSON)
		}
	}
}
