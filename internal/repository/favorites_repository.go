package repository

import (
	"github.com/Aur4ik/AlaRent/internal/config"
	"github.com/Aur4ik/AlaRent/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AddFavorite(userID, apartmentID uint) error {
	return config.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "apartment_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"deleted_at": nil,
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).Create(&models.Favorite{
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
		Joins("JOIN apartments ON apartments.id = favorites.apartment_id AND apartments.deleted_at IS NULL").
		Preload("Apartment").
		Preload("Apartment.Photos").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error
	return favorites, err
}
