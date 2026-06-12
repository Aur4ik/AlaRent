package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return errors.New("invalid DBHost")
	}
	if c.DBName == "" {
		return errors.New("invalid DBName")
	}
		if c.DBUser == "" {
		return errors.New("invalid DBUser")
	}
	if c.DBPassword == "" {
		return errors.New("invalid DBPassword")
	}
		if c.DBPort == "" {
		return errors.New("invalid DBPort")
	}
	if c.JWTSecret == "" {
		return errors.New("invalid JWT")
	}
	return nil
}
