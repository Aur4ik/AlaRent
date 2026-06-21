package repository

import (
	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
)

func CreateAppartaments(appart *models.Apartment) error {
	return config.DB.Create(appart).Error
}
func GetAllApartments() ([]models.Apartment, error) {
	var apartments []models.Apartment

	err := config.DB.Find(&apartments).Error

	return apartments, err
}

func GetApartmentByID(id uint) (*models.Apartment, error) {
	var apartment models.Apartment

	err := config.DB.First(&apartment, id).Error

	return &apartment, err
}

func UpdateApartment(apartment *models.Apartment) error {
	return config.DB.Save(apartment).Error
}

func DeleteApartment(apartment *models.Apartment) error {
	return config.DB.Delete(apartment).Error
}
