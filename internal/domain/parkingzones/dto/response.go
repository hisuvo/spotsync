package dto

import "time"

type ParkingZoneResponse struct {
	ID            uint64  `json:"id"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	TotalCapacity int     `json:"total_capacity"`
	PricePerHour  float64 `json:"price_per_hour"`
	AvailableSpots  int     `json:"available_spots"`
	CreatedAt     time.Time  `json:"created_at"`
}