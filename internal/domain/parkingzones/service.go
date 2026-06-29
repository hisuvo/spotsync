package parkingzones

import (
	"fmt"
	"spotsync/internal/domain/parkingzones/dto"
)

type ReservationCounter interface { // todo
	CountActiveReservations(zoneID uint64) (int64, error)
}

type Service interface {
	Create(req *dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	GetAll() ([]dto.ParkingZoneResponse, error)
	FindResponseByID(id uint64) (*dto.ParkingZoneResponse, error)
	Update(id uint64, req *dto.UpdateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	Delete(id uint64) error
}

type service struct {
	repo Repository
	reservationRepo ReservationCounter //todo
}

func NewParkingZoneService(repo Repository,reservationRepo ReservationCounter) Service {
	return &service{
		repo: repo,
		reservationRepo: reservationRepo,
	}
}

func (s *service) Create(req *dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone := &ParkingZone{
		Name: req.Name,
		Type: req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour: req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, fmt.Errorf("create parking zone: %w", err)
	}

	 response := dto.ParkingZoneResponse{
		ID:            uint64(zone.ID),
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
	}

	return &response, nil
}

func (s *service) GetAll() ([]dto.ParkingZoneResponse, error) {
	return s.repo.GetAll()
}

func (s *service) FindResponseByID(id uint64) (*dto.ParkingZoneResponse, error) {
	// Return repo result directly — it already includes AvailableSpots from the JOIN query.
	return s.repo.FindResponseByID(id)
}

func (s *service) Update(id uint64, req *dto.UpdateParkingZoneRequest) (*dto.ParkingZoneResponse, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	zone.Name = req.Name
	zone.Type = req.Type
	zone.TotalCapacity = req.TotalCapacity
	zone.PricePerHour = req.PricePerHour

	if err := s.repo.Update(zone); err != nil {
		return nil, err
	}

	return &dto.ParkingZoneResponse{
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		
	}, nil
}

func (s *service) Delete(id uint64) error {
	return s.repo.Delete(id)
}