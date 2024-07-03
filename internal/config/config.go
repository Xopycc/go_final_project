package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DB   string
}

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(fmt.Errorf("config - New - Error loading .env file: %w", err))
	}

	port := os.Getenv("TODO_PORT")
	db := os.Getenv("TODO_DB")

	return &Config{
		Port: port,
		DB:   db,
	}, nil
}