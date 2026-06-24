package service

import (
	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
)

func AddFavorite(userID, apartmentID uint) error {
	if _, err := GetApartmentByID(apartmentID); err != nil {
		return err
	}

	return repository.AddFavorite(userID, apartmentID)
}

func RemoveFavorite(userID, apartmentID uint) error {
	return repository.RemoveFavorite(userID, apartmentID)
}

func GetUserFavorites(userID uint) ([]models.Favorite, error) {
	return repository.GetUserFavorites(userID)
}
