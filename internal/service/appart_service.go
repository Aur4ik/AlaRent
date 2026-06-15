package service

import (
	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
)

func CreateAppartament(apart *models.Apartment) error {
	return repository.CreateAppartaments(apart)
}
func GetAllApartments() ([]models.Apartment, error) {
	return repository.GetAllApartments()
}