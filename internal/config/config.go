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

	return &Config{
		PORT:  fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		DbUrl: os.Getenv("GOOSE_DBSTRING"),
	}
}
