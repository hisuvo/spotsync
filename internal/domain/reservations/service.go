package reservations

import (
	"spotsync/internal/domain/reservations/dto"
)

type Service interface {
	Create(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]Reservation, error)
	GetAllReservations() ([]Reservation, error)
	CancelReservation(userID uint, reservationID uint, isAdmin bool) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// ---------------- CREATE ----------------

func (s *service) Create(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {

	res, err := s.repo.Create(userID, req)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
			ID:           res.ID,
			UserID:       res.UserID,
			ZoneID:       res.ZoneID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			CreatedAt:    res.CreatedAt,
			UpdatedAt:    res.UpdatedAt,

	}, nil
}

// ---------------- GET MY ----------------

func (s *service) GetMyReservations(userID uint) ([]Reservation, error) {
	return s.repo.GetMy(userID)
}

// ---------------- GET ALL ----------------

func (s *service) GetAllReservations() ([]Reservation, error) {
	return s.repo.GetAll()
}

// ---------------- CANCEL ----------------

func (s *service) CancelReservation(userID uint, reservationID uint, isAdmin bool) error {
	return s.repo.Cancel(userID, reservationID, isAdmin)
}