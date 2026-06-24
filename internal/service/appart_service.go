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

func CreateApartment(apart *models.Apartment, photoURLs []string) error {
	if apart.Type == "" {
		apart.Type = "apartment"
	}

	if err := repository.CreateAppartaments(apart); err != nil {
		return err
	}

	if len(photoURLs) > 0 {
		return repository.ReplaceApartmentPhotos(apart.ID, photoURLs)
	}

	return nil
}

func GetAllApartments(filter dto.ApartmentFilter) ([]models.Apartment, error) {
	return repository.GetAllApartments(filter)
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
	if req.Type != nil {
		apartment.Type = *req.Type
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
	if req.HasWasher != nil {
		apartment.HasWasher = *req.HasWasher
	}

	if err := repository.UpdateApartment(apartment); err != nil {
		return nil, err
	}
	if req.PhotoURLs != nil {
		if len(*req.PhotoURLs) > 10 {
			return nil, errors.New("maximum 10 photos allowed")
		}
		if err := repository.ReplaceApartmentPhotos(apartment.ID, *req.PhotoURLs); err != nil {
			return nil, err
		}
	}

	return GetApartmentByID(apartment.ID)
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
