package models

import "gorm.io/gorm"

type Apartment struct {
	gorm.Model

	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Type        string `gorm:"not null;default:'apartment'"`

	Price int `gorm:"not null"`

	District string `gorm:"not null"`
	Address  string `gorm:"not null"`

	Rooms int `gorm:"not null"`
	Floor int `gorm:"not null"`

	HasFurniture bool `gorm:"default:false"`
	HasWifi      bool `gorm:"default:false"`
	HasWasher    bool `gorm:"default:false"`

	OwnerID uint `gorm:"not null"`

	Owner  User             `gorm:"foreignKey:OwnerID"`
	Photos []ApartmentPhoto `gorm:"foreignKey:ApartmentID"`
}

type ApartmentPhoto struct {
	gorm.Model
	ApartmentID uint   `gorm:"not null;index"`
	URL         string `gorm:"not null"`
	IsMain      bool   `gorm:"default:false"`
}

type Favorite struct {
	gorm.Model
	UserID      uint      `gorm:"not null;uniqueIndex:idx_user_apartment"`
	ApartmentID uint      `gorm:"not null;uniqueIndex:idx_user_apartment"`
	User        User      `gorm:"foreignKey:UserID"`
	Apartment   Apartment `gorm:"foreignKey:ApartmentID"`
}
