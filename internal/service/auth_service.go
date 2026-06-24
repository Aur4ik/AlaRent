package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/Aur4ik/AlaRent/internal/dto"
	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
	"github.com/Aur4ik/AlaRent/internal/utils"
)

func Register(user *models.User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	if user.Role == "" {
		user.Role = models.RoleTenant
	}

	if user.Role != models.RoleTenant && user.Role != models.RoleLandlord {
		return errors.New("role must be tenant or landlord")
	}

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

func LoginWithTokens(email, password string) (*models.User, string, string, error) {
	user, err := Login(email, password)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := CreateTokenPair(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func CreateTokenPair(userID uint) (string, string, error) {
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := repository.CreateRefreshToken(&models.RefreshToken{
		UserID: userID,
		Token:  refreshToken,
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func RefreshAccessToken(refreshToken string) (string, error) {
	token, err := repository.GetRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	return utils.GenerateAccessToken(token.UserID)
}

func Logout(refreshToken string) error {
	return repository.DeleteRefreshToken(refreshToken)
}

func GetUserByID(id uint) (*models.User, error) {
	return repository.GetUserByID(id)
}

func UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*models.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.AvatarURL != nil {
		user.AvatarURL = *req.AvatarURL
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	if err := repository.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
