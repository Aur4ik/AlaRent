package service

import (
	"errors"

	"github.com/Aur4ik/AlaRent/internal/dto"
	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrApartmentNotFound  = errors.New("apartment not found")
	ErrApartmentForbidden = errors.New("you can manage only your own apartment")
)

func CreateAppartament(apart *models.Apartment) error {
	return repository.CreateAppartaments(apart)
}

func GetAllApartments() ([]models.Apartment, error) {
	return repository.GetAllApartments()
}

func GetApartmentByID(id uint) (*models.Apartment, error) {
	apartment, err := repository.GetApartmentByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrApartmentNotFound
	}
	return apartment, err
}

func UpdateApartment(id, userID uint, req dto.UpdateApartmentRequest) (*models.Apartment, error) {
	apartment, err := GetApartmentByID(id)
	if err != nil {
		return nil, err
	}

	if apartment.OwnerID != userID {
		return nil, ErrApartmentForbidden
	}

	if req.Title != nil {
		apartment.Title = *req.Title
	}
	if req.Description != nil {
		apartment.Description = *req.Description
	}
	if req.Price != nil {
		apartment.Price = *req.Price
	}
	if req.District != nil {
		apartment.District = *req.District
	}
	if req.Address != nil {
		apartment.Address = *req.Address
	}
	if req.Rooms != nil {
		apartment.Rooms = *req.Rooms
	}
	if req.Floor != nil {
		apartment.Floor = *req.Floor
	}
	if req.HasFurniture != nil {
		apartment.HasFurniture = *req.HasFurniture
	}
	if req.HasWifi != nil {
		apartment.HasWifi = *req.HasWifi
	}

	if err := repository.UpdateApartment(apartment); err != nil {
		return nil, err
	}

	return apartment, nil
}

func DeleteApartment(id, userID uint) error {
	apartment, err := GetApartmentByID(id)
	if err != nil {
		return err
	}

	if apartment.OwnerID != userID {
		return ErrApartmentForbidden
	}

	return repository.DeleteApartment(apartment)
}
