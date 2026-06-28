package database

import (
	"spotsync/internal/domain/user"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	 db.AutoMigrate(&user.User{})
}