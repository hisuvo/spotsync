package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseUrl string
}

func LoadEnv() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found", err)
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		DatabaseUrl: os.Getenv("DB_URL"),
	}
}