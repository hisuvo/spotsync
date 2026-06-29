package database

import (
	"spotsync/internal/domain/user"
	"spotsync/internal/domain/parkingzones"
	"spotsync/internal/domain/reservations"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&parkingzones.ParkingZone{})
	db.AutoMigrate(&reservations.Reservation{})
}