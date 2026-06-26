package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Aur4ik/AlaRent/internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = buildPostgresDSN()
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Bug fix: original error was swallowed — now we log it for easier debugging
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
		&models.Apartment{},
		&models.ApartmentPhoto{},
		&models.Favorite{},
		&models.Conversation{},
		&models.Message{},
	); err != nil {
		log.Fatalf("Auto migration failed: %v", err)
	}

	DB = db
	log.Println("Database connected and migrations applied")
}

func buildPostgresDSN() string {
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		sslMode,
	)
}
