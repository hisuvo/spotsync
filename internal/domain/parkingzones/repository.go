package parkingzones

import (
	"errors"

	"gorm.io/gorm"
)

var ErrParkingZoneNotFound = errors.New("parking zone not found")
var ErrParkingZoneAlreadyExists = errors.New("parking zone already exists")

type Repository interface {
	Create(parkingZone *ParkingZone)  error
	GetAll() (*[]ParkingZone, error)
	FindById(id uint64) (*ParkingZone, error)
	Update(parkingZone *ParkingZone)  error
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

func (r *repository) GetAll() (*[]ParkingZone, error) {
	var zones *[]ParkingZone

	err := r.db.Find(&zones).Error

	if err != nil {
		return nil, ErrParkingZoneNotFound
	}
	return zones, nil
}

func (r *repository) FindById(id uint64) (*ParkingZone, error) {
	var zone ParkingZone
	if err := r.db.Where("id = ?", id).First(&zone).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

// func (r *repository) FindByID(id uint64) (*dto.ParkingZoneResponse, error) {
// 	var zone dto.ParkingZoneResponse

// 	subQuery := r.db.
// 		Model(&Reservation{}).
// 		Select("COUNT(*)").
// 		Where("reservations.parking_zone_id = parking_zones.id").
// 		Where("status = ?", "active")

// 	err := r.db.
// 		Model(&ParkingZone{}).
// 		Select(`
// 			id,
// 			name,
// 			type,
// 			total_capacity,
// 			price_per_hour,
// 			created_at,
// 			total_capacity - (?) AS available_spots
// 		`, subQuery).
// 		Where("id = ?", id).
// 		Scan(&zone).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &zone, nil
// }

func (r *repository) Update(zone *ParkingZone)  error {
	if err := r.db.Save(zone).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id uint64) error {
	if err := r.db.Where("id = ?", id).Delete(&ParkingZone{}).Error; err != nil {
		return err
	}
	return nil
}