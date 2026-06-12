package models

import(
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Phone     string    `gorm:"not null"` 
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}