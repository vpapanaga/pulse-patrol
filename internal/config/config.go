package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application settings retrieved from environment variables.
type Config struct {
	RestPort string
	GrpcPort string
	DBUrl    string
}

// LoadConfig initializes the configuration by reading the .env file.
// If the .env file is missing, it logs a warning and continues using system environment variables.
func LoadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found. Falling back to system environment variables.")
	}
}

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value if present, otherwise it returns the provided fallback string.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
