package models

import "gorm.io/gorm"

type Apartment struct {
	gorm.Model

	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`

	Price int `gorm:"not null"`

	District string `gorm:"not null"`
	Address  string `gorm:"not null"`

	Rooms int `gorm:"not null"`
	Floor int `gorm:"not null"`

	HasFurniture bool `gorm:"default:false"`
	HasWifi      bool `gorm:"default:false"`

	OwnerID uint `gorm:"not null"`

	Owner User `gorm:"foreignKey:OwnerID"`
}