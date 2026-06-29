package database

import (
	"spotsync/internal/domain/user"

	"spotsync/internal/domain/parkingzones"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&parkingzones.ParkingZone{})
}