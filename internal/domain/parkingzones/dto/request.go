package dto

type CreateParkingZoneRequest struct {
	Name          string  `json:"name" validate:"required,min=3,max=255"`
	Type          string  `json:"type" validate:"required,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"required,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"required,gt=0"`
}

type UpdateParkingZoneRequest struct {
	Name          string  `json:"name" validate:"omitempty,min=3,max=255"`
	Type          string  `json:"type" validate:"omitempty,oneof=general ev_charging covered"`
	TotalCapacity int     `json:"total_capacity" validate:"omitempty,gt=0"`
	PricePerHour  float64 `json:"price_per_hour" validate:"omitempty,gt=0"`
}