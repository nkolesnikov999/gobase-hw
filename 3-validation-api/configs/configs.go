package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// SMTPConfig contains SMTP credentials and server address.
type SMTPConfig struct {
	Email    string
	Password string
	Address  string
}

// Config is the root application configuration.
type Config struct {
	SMTP SMTPConfig
}

// LoadConfig loads configuration from environment variables with sane defaults.
//
// Environment variables:
// - SMTP_EMAIL
// - SMTP_PASSWORD
// - SMTP_ADDRESS (host:port)
func LoadConfig() Config {
	// Load .env if present. Missing file is not an error.
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}
	return Config{
		SMTP: SMTPConfig{
			Email:    getEnv("SMTP_EMAIL", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			Address:  getEnv("SMTP_ADDRESS", "localhost:1025"),
		},
	}
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
