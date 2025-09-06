package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DbUrl       string
	TenorAPIKey string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// first try to get the port from the environment variable (for Cloud Run)
	port := os.Getenv("PORT")
	var env string
	if port != "" {
		env = "production"
	}

	// fallback to APP_PORT for local development
	if port == "" {
		port = os.Getenv("APP_PORT")
		env = "development"
	}

	// fallback to 8080 if no port is set
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:        fmt.Sprintf(":%s", port),
		Env:         env,
		DbUrl:       os.Getenv("GOOSE_DBSTRING"),
		TenorAPIKey: os.Getenv("TENOR_API_KEY"),
	}
}
