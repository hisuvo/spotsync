package dto

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required,gt=0"`
	LicensePlate string `gorm:"type:varchar(15);not null" json:"license_plate" validate:"required,max=15"`
}