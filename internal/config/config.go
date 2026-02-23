// Package config handles environment variable loading and application configuration.
// It supports both local .env files and system environment variables for Docker/Cloud-native deployments.
package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the application settings.
type Config struct {
	RestPort   string
	GrpcPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// GlobalConfig stores the loaded settings to be accessed across the application.
var GlobalConfig *Config

// LoadConfig initializes the configuration.
// It attempts to load a .env file but ignores failures if variables are provided via system environment (Docker).
func LoadConfig() {
	// Attempt to load the .env file. We use the blank identifier (_) to ignore the error
	// because in Docker Compose/Production, variables are injected directly into the OS environment.
	_ = godotenv.Load(".env")

	GlobalConfig = &Config{
		RestPort:   GetEnv("REST_PORT", "8080"),
		GrpcPort:   GetEnv("GRPC_PORT", "50051"),
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     GetEnv("DB_PORT", "5432"),
		DBUser:     GetEnv("DB_USER", "postgres"),
		DBPassword: GetEnv("DB_PASSWORD", "password"),
		DBName:     GetEnv("DB_NAME", "pulse_patrol_db"),
	}

	log.Printf("Configuration initialized. REST Port: %s, DB Host: %s", GlobalConfig.RestPort, GlobalConfig.DBHost)
}

// GetDatabaseDSN constructs the connection string for the pgx driver.
// It follows the format: postgres://user:password@host:port/dbname?sslmode=disable
func GetDatabaseDSN() string {
	if GlobalConfig == nil {
		LoadConfig()
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		GlobalConfig.DBUser,
		GlobalConfig.DBPassword,
		GlobalConfig.DBHost,
		GlobalConfig.DBPort,
		GlobalConfig.DBName,
	)
}

// GetEnv is a helper function to read an environment variable or return a fallback value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
