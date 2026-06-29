package reservations

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	gorm.Model
	UserID    uint64
	SpotID    uint64
	StartTime time.Time
	EndTime   time.Time
	Status    string
	
}