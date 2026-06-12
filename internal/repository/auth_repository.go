package repository

import (
	"github.com/Aur4ik/AlaRent/pkg/config"
	"github.com/Aur4ik/AlaRent/internal/models"
)

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserRole(id uint, role string) error {
	return config.DB.
		Model(&models.User{}).
		Where("id = ?", id).
		Update("role", role).Error
}

func CountAdmins() (int64, error) {
	var count int64
	err := config.DB.
		Model(&models.User{}).
		Where("role = ?", "admin").
		Count(&count).Error
	return count, err
}
