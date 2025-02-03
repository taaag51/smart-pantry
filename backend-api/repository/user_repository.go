package repository

import (
	"github.com/taaag51/smart-pantry/backend-api/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	CreateUser(user *model.User) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetUserByID(id uint) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(user *model.User) (model.User, error) {
	if err := ur.db.Create(user).Error; err != nil {
		return model.User{}, err
	}
	return *user, nil
}

func (ur *userRepository) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ur *userRepository) GetUserByID(id uint) (model.User, error) {
	var user model.User
	if err := ur.db.First(&user, id).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
