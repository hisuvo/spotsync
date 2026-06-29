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
	// Only load .env in local dev (safe for production)
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // safe default
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	env := os.Getenv("APP_ENV")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if env == "development" {
			jwtSecret = "dev-secret-key"
		} else {
			log.Fatal("JWT_SECRET is not set")
		}
	}

	duration := 24 * time.Hour

	if value := os.Getenv("TOKEN_DURATION"); value != "" {
		d, err := time.ParseDuration(value)
		if err != nil {
			log.Fatalf("Invalid TOKEN_DURATION: %v", err)
		}
		duration = d
	}

	return &Config{
		Port:          port,
		DatabaseURL:   dbURL,
		JWTSecret:     jwtSecret,
		TokenDuration: duration,
	}
}