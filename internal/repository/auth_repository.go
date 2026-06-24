package repository

import (
	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
)

func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

func CreateRefreshToken(token *models.RefreshToken) error {
	return config.DB.Create(token).Error
}

func GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := config.DB.Where("token = ?", token).First(&refreshToken).Error
	return &refreshToken, err
}

func DeleteRefreshToken(token string) error {
	return config.DB.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}
