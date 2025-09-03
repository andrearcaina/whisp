package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT  string
	DbUrl string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT") // for Cloud Run
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080" // Default fallback
	}

	return &Config{
		PORT:  fmt.Sprintf(":%s", port),
		DbUrl: os.Getenv("GOOSE_DBSTRING"),
	}
}
