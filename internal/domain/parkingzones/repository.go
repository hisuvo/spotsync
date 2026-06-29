package parkingzones

import (
	"errors"
	"spotsync/internal/domain/parkingzones/dto"

	"gorm.io/gorm"
)

var ErrParkingZoneNotFound = errors.New("parking zone not found")
var ErrParkingZoneAlreadyExists = errors.New("parking zone already exists")

type Repository interface {
	Create(parkingZone *ParkingZone) error
	GetAll() ([]dto.ParkingZoneResponse, error)
	FindByID(id uint64) (*ParkingZone, error)
	FindResponseByID(id uint64) (*dto.ParkingZoneResponse, error)
	Update(id uint64,parkingZone *ParkingZone) error
	Delete(id uint64) error
}

type repository struct {
	db *gorm.DB
}

func NewParkingZoneRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(zone *ParkingZone)  error {
	if err := r.db.Create(zone).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAll() ([]dto.ParkingZoneResponse, error) {
	var zones []dto.ParkingZoneResponse

	err := r.db.
		Table("parking_zones").
		Select(`
			parking_zones.id,
			parking_zones.name,
			parking_zones.type,
			parking_zones.total_capacity,
			parking_zones.price_per_hour,
			parking_zones.created_at,
			parking_zones.total_capacity - COALESCE(COUNT(r.id), 0) AS available_spots
		`).
		Joins(`LEFT JOIN reservations r ON r.zone_id = parking_zones.id AND r.status = 'active' AND r.deleted_at IS NULL`).
		Where("parking_zones.deleted_at IS NULL").
		Group("parking_zones.id").
		Scan(&zones).Error

	if err != nil {
		return nil, err
	}
	return zones, nil
}

func (r *repository) FindByID(id uint64) (*ParkingZone, error) {
	var zone ParkingZone
	if err := r.db.Where("id = ?", id).First(&zone).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *repository) FindResponseByID(id uint64) (*dto.ParkingZoneResponse, error) {
	var zone dto.ParkingZoneResponse

	err := r.db.
		Table("parking_zones").
		Select(`
			parking_zones.id,
			parking_zones.name,
			parking_zones.type,
			parking_zones.total_capacity,
			parking_zones.price_per_hour,
			parking_zones.created_at,
			parking_zones.total_capacity - COALESCE(COUNT(r.id), 0) AS available_spots
		`).
		Joins(`LEFT JOIN reservations r ON r.zone_id = parking_zones.id AND r.status = 'active' AND r.deleted_at IS NULL`).
		Where("parking_zones.id = ?", id).
		Where("parking_zones.deleted_at IS NULL").
		Group("parking_zones.id").
		Scan(&zone).Error

	if err != nil {
		return nil, err
	}
	if zone.ID == 0 {
		return nil, ErrParkingZoneNotFound
	}

	return &zone, nil
}

func (r *repository) Update(id uint64, zone *ParkingZone) error {
	return r.db.Model(&ParkingZone{}).
		Where("id = ?", id).
		Updates(zone).Error
}

func (r *repository) Delete(id uint64) error {
	if err := r.db.Where("id = ?", id).Delete(&ParkingZone{}).Error; err != nil {
		return err
	}
	return nil
}