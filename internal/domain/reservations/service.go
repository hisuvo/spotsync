package reservations

import (
	"spotsync/internal/domain/reservations/dto"
)

type Service interface {
	Create(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	GetAllReservations() ([]dto.AdminReservationResponse, error)
	CancelReservation(userID uint, reservationID uint, isAdmin bool) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}


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

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	resvs, err := s.repo.GetMy(userID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.MyReservationResponse, 0, len(resvs))
	for _, r := range resvs {
		resp = append(resp, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.MyReservationZoneDTO{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt:    r.CreatedAt,
		})
	}
	return resp, nil
}

func (s *service) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	resvs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.AdminReservationResponse, 0, len(resvs))
	for _, r := range resvs {
		resp = append(resp, dto.AdminReservationResponse{
			ID:           r.ID,
			UserID:       r.UserID,
			User: dto.AdminUserDTO{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
				Role:  r.User.Role,
			},
			ZoneID:       r.ZoneID,
			Zone: dto.AdminZoneDTO{
				ID:            r.Zone.ID,
				Name:          r.Zone.Name,
				Type:          r.Zone.Type,
				TotalCapacity: r.Zone.TotalCapacity,
				PricePerHour:  r.Zone.PricePerHour,
			},
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}
	return resp, nil
}

func (s *service) CancelReservation(userID uint, reservationID uint, isAdmin bool) error {
	return s.repo.Cancel(userID, reservationID, isAdmin)
}