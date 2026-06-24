package repository

import (
	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
)

func AddFavorite(userID, apartmentID uint) error {
	return config.DB.FirstOrCreate(&models.Favorite{}, models.Favorite{
		UserID:      userID,
		ApartmentID: apartmentID,
	}).Error
}

func RemoveFavorite(userID, apartmentID uint) error {
	return config.DB.Where("user_id = ? AND apartment_id = ?", userID, apartmentID).
		Delete(&models.Favorite{}).Error
}

func GetUserFavorites(userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := config.DB.
		Preload("Apartment").
		Preload("Apartment.Photos").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error
	return favorites, err
}
