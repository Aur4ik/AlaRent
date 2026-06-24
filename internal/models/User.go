package models

import "gorm.io/gorm"

const (
	RoleTenant   = "tenant"
	RoleLandlord = "landlord"
)

type User struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Phone     string `gorm:"not null"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null;default:'tenant'"`
	AvatarURL string
	Bio       string `gorm:"type:text"`
}

type RefreshToken struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Token  string `gorm:"uniqueIndex;not null"`
	User   User   `gorm:"foreignKey:UserID"`
}
