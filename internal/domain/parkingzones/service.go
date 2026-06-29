package parkingzones

import (
	"fmt"
	"spotsync/internal/domain/parkingzones/dto"
)

type Service interface {
	Create(req *dto.CreateParkingZoneRequest) (*dto.ParkingZoneResponse, error)
	GetAll() (*[]dto.ParkingZoneResponse, error)
	FindById(id uint64) (*dto.ParkingZoneResponse, error)
	Update(id uint64, req *dto.UpdateParkingZoneRequest) (*dto.UpdateParkingZoneRequest, error)
	Delete(id uint64) error
}

type service struct {
	repo Repository
}

func NewParkingZoneService(repo Repository) Service {
	return &service{
		repo: repo,
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

func (s *service) GetAll()(*[]dto.ParkingZoneResponse,error){
	zones, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var responses []dto.ParkingZoneResponse
	for _, zone := range *zones {
		responses = append(responses, dto.ParkingZoneResponse{
			ID:            uint64(zone.ID),
			Name:          zone.Name,
			Type:          zone.Type,
			TotalCapacity: zone.TotalCapacity,
			PricePerHour:  zone.PricePerHour,
			CreatedAt:     zone.CreatedAt,
		})
	}

	return &responses, nil
}

func (s *service) FindById(id uint64) (*dto.ParkingZoneResponse, error) {
	zone, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	// respnse := zone.ToParkingZoneResponse()

	return &dto.ParkingZoneResponse{
		ID:            uint64(zone.ID),
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
	}, nil
}

func (s *service) Update(id uint64, req *dto.UpdateParkingZoneRequest) (*dto.UpdateParkingZoneRequest, error) {
	zone, err := s.repo.FindById(id)
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

	return &dto.UpdateParkingZoneRequest{
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		
	}, nil
}

func (s *service) Delete(id uint64) error {
	return s.repo.Delete(id)
}