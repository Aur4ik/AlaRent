package repository

import (
	"strings"

	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/dto"
	"github.com/Aur4ik/AlaRent/internal/models"
	"gorm.io/gorm"
)

func CreateAppartaments(appart *models.Apartment) error {
	return config.DB.Create(appart).Error
}

func FindDuplicateApartment(ownerID uint, apartmentID uint, title, address, apartmentType string, rooms, floor int) (*models.Apartment, error) {
	var apartment models.Apartment

	query := config.DB.
		Where("owner_id = ?", ownerID).
		Where("lower(title) = lower(?)", title).
		Where("lower(address) = lower(?)", address).
		Where("type = ?", apartmentType).
		Where("rooms = ? AND floor = ?", rooms, floor)

	if apartmentID != 0 {
		query = query.Where("id <> ?", apartmentID)
	}

	err := query.First(&apartment).Error
	return &apartment, err
}

func GetAllApartments(filter dto.ApartmentFilter) ([]models.Apartment, error) {
	var apartments []models.Apartment

	query := config.DB.Model(&models.Apartment{}).Preload("Photos").Preload("Owner")

	if filter.Query != "" {
		search := "%" + strings.ToLower(filter.Query) + "%"
		query = query.Where("lower(title) LIKE ? OR lower(description) LIKE ? OR lower(address) LIKE ?", search, search, search)
	}
	if filter.District != "" {
		query = query.Where("district = ?", filter.District)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}
	if filter.Rooms > 0 {
		query = query.Where("rooms = ?", filter.Rooms)
	}
	if filter.HasFurniture != nil {
		query = query.Where("has_furniture = ?", *filter.HasFurniture)
	}
	if filter.HasWifi != nil {
		query = query.Where("has_wifi = ?", *filter.HasWifi)
	}
	if filter.HasWasher != nil {
		query = query.Where("has_washer = ?", *filter.HasWasher)
	}

	switch filter.Sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "oldest":
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Find(&apartments).Error

	return apartments, err
}

func GetApartmentByID(id uint) (*models.Apartment, error) {
	var apartment models.Apartment

	err := config.DB.Preload("Photos").Preload("Owner").First(&apartment, id).Error

	return &apartment, err
}

func UpdateApartment(apartment *models.Apartment) error {
	return config.DB.Save(apartment).Error
}

func DeleteApartment(apartment *models.Apartment) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("apartment_id = ?", apartment.ID).Delete(&models.Favorite{}).Error; err != nil {
			return err
		}
		if err := tx.Where("apartment_id = ?", apartment.ID).Delete(&models.ApartmentPhoto{}).Error; err != nil {
			return err
		}

		var conversations []models.Conversation
		if err := tx.Where("apartment_id = ?", apartment.ID).Find(&conversations).Error; err != nil {
			return err
		}
		for _, conversation := range conversations {
			if err := tx.Where("conversation_id = ?", conversation.ID).Delete(&models.Message{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("apartment_id = ?", apartment.ID).Delete(&models.Conversation{}).Error; err != nil {
			return err
		}

		return tx.Delete(apartment).Error
	})
}

func ReplaceApartmentPhotos(apartmentID uint, photoURLs []string) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("apartment_id = ?", apartmentID).Delete(&models.ApartmentPhoto{}).Error; err != nil {
			return err
		}

		for i, url := range photoURLs {
			photo := models.ApartmentPhoto{
				ApartmentID: apartmentID,
				URL:         url,
				IsMain:      i == 0,
			}
			if err := tx.Create(&photo).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
