package database

import (
	"fmt"
	"spotsync/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	connStr := cfg.DatabaseURL

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		TranslateError: true,
		// Show only warnings
		Logger: logger.Default.LogMode(logger.Warn),
		// Show everything (development):
		// Logger: logger.Default.LogMode(logger.Info), 
	})

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil
	}

	Migrate(db) // this auto create table in DB
	fmt.Println("Successfully database connected done!")
	return db
}