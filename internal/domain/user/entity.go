package user

import (
	"spotsync/internal/domain/user/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User Roles
const (
	RoleAdmin  = "admin"
	RoleDriver = "driver"
)

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(50);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `json:"password" gorm:"type:varchar(250);not null"`
	Role     string `json:"role" gorm:"type:varchar(20);default:driver;not null"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ToResponse() *dto.UserResponse{
	return &dto.UserResponse{
		ID:       uint64(u.ID),
		Name:     u.Name,
		Email:    u.Email,
		Role:     u.Role,
		CreateAt: u.CreatedAt.Unix(),
		UpdateAt: u.UpdatedAt.Unix(),
	}
}