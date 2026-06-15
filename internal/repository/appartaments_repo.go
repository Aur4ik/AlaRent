package repository

import (
	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
)

func CreateAppartaments(appart *models.Apartment) error{
	return config.DB.Create(appart).Error
}
func GetAllApartments() ([]models.Apartment, error){
	var apartments []models.Apartment

	err := config.DB.Find(&apartments).Error

	return apartments, err
}