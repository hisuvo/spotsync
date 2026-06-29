package reservations

import (
	"errors"

	"spotsync/internal/domain/parkingzones"
	"spotsync/internal/domain/reservations/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull            = errors.New("parking zone is full")
	ErrReservationNotFound = errors.New("reservation not found")
	ErrUnauthorizedCancel  = errors.New("unauthorized to cancel this reservation")
	ErrAlreadyProcessed    = errors.New("reservation already completed or cancelled")
)

type Repository interface {
	Create(userID uint, req *dto.CreateReservationRequest) (*Reservation, error)
	CountActiveReservations(zoneID uint64) (int64, error) // todo
	GetMy(userID uint) ([]Reservation, error)
	GetAll() ([]Reservation, error)
	Cancel(userID uint, reservationID uint, isAdmin bool) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// ---------------- CREATE (SAFE TRANSACTION + LOCK) ----------------

func (r *repository) Create(userID uint, req *dto.CreateReservationRequest) (*Reservation, error) {

	var reservation Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {

		var zone parkingzones.ParkingZone

		// Lock row (prevents race condition)
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&zone, req.ZoneID).Error; err != nil {
			return err
		}

		var count int64

		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ?", req.ZoneID, StatusActive).
			Count(&count).Error; err != nil {
			return err
		}

		if count >= int64(zone.TotalCapacity) {
			return ErrZoneFull
		}

		reservation = Reservation{
			UserID:       userID,
			ZoneID:       req.ZoneID,
			LicensePlate: req.LicensePlate,
			Status:       string(StatusActive),
		}

		return tx.Create(&reservation).Error
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *repository) CountActiveReservations(zoneID uint64) (int64, error) {
	var count int64

	err := r.db.
		Table("reservations").
		Where("zone_id = ?", zoneID).
		Where("status = ?", "active").
		Count(&count).Error

	return count, err
}
// ---------------- GET MY ----------------

func (r *repository) GetMy(userID uint) ([]Reservation, error) {
	var res []Reservation

	err := r.db.
		Preload("Zone").
		Where("user_id = ?", userID).
		Find(&res).Error

	return res, err
}

// ---------------- GET ALL ----------------

func (r *repository) GetAll() ([]Reservation, error) {
	var res []Reservation

	err := r.db.
		Preload("User").
		Preload("Zone").
		Find(&res).Error

	return res, err
}

// ---------------- CANCEL ----------------

func (r *repository) Cancel(userID uint, reservationID uint, isAdmin bool) error {

	var resv Reservation

	err := r.db.First(&resv, reservationID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReservationNotFound
		}
		return err
	}

	if resv.UserID != userID && !isAdmin {
		return ErrUnauthorizedCancel
	}

	if resv.Status == string(StatusCompleted) || resv.Status == string(StatusCancelled) {
		return ErrAlreadyProcessed
	}

	resv.Status = string(StatusCancelled)

	return r.db.Save(&resv).Error
}