package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	TokenDuration time.Duration
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	env := os.Getenv("APP_ENV")

	if jwtSecret == "" {
		if env == "development" {
			jwtSecret = "dev-secret-key"
		} else {
			log.Fatal("JWT_SECRET is not set")
		}
	}

	duration := 24 * time.Hour

	if value := os.Getenv("TOKEN_DURATION"); value != "" {
		var err error
		duration, err = time.ParseDuration(value)
		if err != nil {
			log.Fatalf("Invalid TOKEN_DURATION: %v", err)
		}
	}

	return &Config{
		Port:          port,
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		TokenDuration: duration,
	}
}