package user

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrorAlreadyExist = errors.New("user with this email already exist")
	ErrorUserNotFound = errors.New("user not found")
)


type UserRepository interface {
	CreateUser(user *User) error
	FindUserByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) UserRepository {
	return &repository{db}
}

func (r *repository) CreateUser(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrorAlreadyExist
		}
		return result.Error	
	}
	return nil
}

func (r *repository) FindUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorUserNotFound
		}
		return nil, err
	}
	return &user, nil
}