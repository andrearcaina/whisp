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
	KlipyAPIKey string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// first try to get the port from the environment variable
	port := os.Getenv("PORT")
	var env string
	var dbUrl string
	if port != "" {
		env = "production"
		dbUrl = os.Getenv("PROD_DBSTRING")
	}

	// fallback to APP_PORT for local development
	if port == "" {
		port = os.Getenv("APP_PORT")
		env = "development"
		dbUrl = os.Getenv("DEV_DBSTRING")
	}

	// fallback to 8080 if no port is set
	if port == "" {
		port = "8080"
		dbUrl = os.Getenv("DEV_DBSTRING")
	}

	return &Config{
		Port:        fmt.Sprintf(":%s", port),
		Env:         env,
		DbUrl:       dbUrl,
		TenorAPIKey: os.Getenv("TENOR_API_KEY"),
		KlipyAPIKey: os.Getenv("KLIPY_API_KEY"),
	}
}
