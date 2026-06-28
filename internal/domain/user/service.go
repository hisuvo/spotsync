package user

import (
	"errors"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"
	"strconv"
)

var ErrInvalideCredentials = errors.New("User email not found")

type UserService interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	LoginUser(email, password string) (*dto.LoginResponse, error)
}


type service struct {
	userRepo UserRepository
	jwtService auth.JWTService
}

func NewService(userRepo UserRepository, jwtService auth.JWTService) UserService {
	return &service{userRepo, jwtService}
}

func (s *service) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {

	user :=User{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
		Role: req.Role,
	}
	
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	if user.Role == "" {
		user.Role = "driver"
	}

	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, err
	}

	// token, err := s.jwtService.GenerateToken(strconv.FormatUint(uint64(user.ID), 10), user.Name, user.Email)
	// if err != nil {
	// 	return nil, err
	// }

	response := dto.UserResponse{
		ID:       uint64(user.ID),
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		CreateAt: user.CreatedAt.Unix(),
		UpdateAt: user.UpdatedAt.Unix(),
	}

	return &response, nil
}

func (s *service) LoginUser(email, password string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindUserByEmail(email)

	if err != nil {
		return nil, ErrInvalideCredentials
	}

	if err := user.CheckPassword(password); err != nil {
		return nil, ErrInvalideCredentials
	}

	token, err := s.jwtService.GenerateToken(strconv.FormatUint(uint64(user.ID), 10), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	response := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:       uint64(user.ID),
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			CreateAt: user.CreatedAt.Unix(),
			UpdateAt: user.UpdatedAt.Unix(),
		},
	}

	return &response, nil
}