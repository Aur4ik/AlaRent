package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
)

func Register(user *models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if user.Role != "" && user.Role != models.RoleTenant {
		return errors.New("only tenant role is allowed on registration")
	}
	user.Role = models.RoleTenant

	return repository.CreateUser(user)
}
func Login(email, password string) (*models.User, error) {
	user, err := repository.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	return repository.GetUserByID(id)
}
