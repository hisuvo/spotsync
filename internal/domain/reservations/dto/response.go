package dto

import "time"

type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MyReservationZoneDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type MyReservationResponse struct {
	ID           uint                 `json:"id"`
	LicensePlate string               `json:"license_plate"`
	Status       string               `json:"status"`
	Zone         MyReservationZoneDTO `json:"zone"`
	CreatedAt    time.Time            `json:"created_at"`
}

type AdminUserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type AdminZoneDTO struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	TotalCapacity int     `json:"total_capacity"`
	PricePerHour  float64 `json:"price_per_hour"`
}

type AdminReservationResponse struct {
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	User         AdminUserDTO `json:"user"`
	ZoneID       uint         `json:"zone_id"`
	Zone         AdminZoneDTO `json:"zone"`
	LicensePlate string       `json:"license_plate"`
	Status       string       `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}