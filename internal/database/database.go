package database

import (
	"fmt"
	"spotsync/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	connStr := cfg.DatabaseURL

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}

	fmt.Println("Successfully database connected done!")
	return db
}