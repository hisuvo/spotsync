package reservations

import (
	"spotsync/internal/domain/parkingzones"
	"spotsync/internal/domain/reservations/dto"
	"spotsync/internal/domain/user"

	"gorm.io/gorm"
)

type ReservationStatus string

const (
	StatusActive    ReservationStatus = "active"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
)

type Reservation struct {
	gorm.Model
	UserID uint      `gorm:"not null;index" json:"user_id"`
	User   user.User `gorm:"foreignKey:UserID"`
	// Composite index on (zone_id, status) so COUNT(*) WHERE zone_id=? AND status='active'
	// uses an index scan instead of a full table scan — fixes the 500ms query latency.
	ZoneID       uint                     `gorm:"not null;index:idx_zone_status" json:"zone_id"`
	Zone         parkingzones.ParkingZone `gorm:"foreignKey:ZoneID"`
	LicensePlate string                   `gorm:"type:varchar(15);not null" json:"license_plate"`
	Status       string                   `gorm:"type:varchar(20);not null;default:active;index:idx_zone_status" json:"status"`
}

func (resv *Reservation) ToResponse() *dto.ReservationResponse {
	return &dto.ReservationResponse{
		ID:           uint(resv.ID),
		UserID:       resv.UserID,
		ZoneID:       resv.ZoneID,
		LicensePlate: resv.LicensePlate,
		Status:       resv.Status,
		CreatedAt:    resv.CreatedAt,
		UpdatedAt:    resv.UpdatedAt,
	}
}