package parkingzones

import (
	"spotsync/internal/domain/parkingzones/dto"

	"gorm.io/gorm"
)

type ParkingZoneType string

const (
	ZoneGeneral    ParkingZoneType = "general"
	ZoneEVCharging ParkingZoneType = "ev_charging"
	ZoneCovered    ParkingZoneType = "covered"
)

type ParkingZone struct {
	gorm.Model
	Name          string  `json:"name" gorm:"type:varchar(255);not null"`
	Type          string  `json:"type" gorm:"type:varchar(50);not null"`
	TotalCapacity int     `json:"total_capacity" gorm:"type:int;check:total_capacity > 0;not null"`
	PricePerHour  float64 `json:"price_per_hour" gorm:"type:decimal(10,2);check:price_per_hour > 0;not null"`
}

func(pz *ParkingZone) ToParkingZoneResponse() *dto.ParkingZoneResponse {
	return &dto.ParkingZoneResponse{
		ID:            uint64(pz.ID),
		Name:          pz.Name,
		Type:          pz.Type,
		TotalCapacity: pz.TotalCapacity,
		PricePerHour:  pz.PricePerHour,
		CreatedAt:     pz.CreatedAt,
	}
}