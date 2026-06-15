package dto

import "github.com/Aur4ik/AlaRent/internal/models"

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"    binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

type CreateApartmentRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	Price        int    `json:"price" binding:"required,gt=0"`
	District     string `json:"district" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Rooms        int    `json:"rooms" binding:"required,gt=0"`
	Floor        int    `json:"floor" binding:"required,gte=0"`
	HasFurniture bool   `json:"has_furniture"`
	HasWifi      bool   `json:"has_wifi"`
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Role:  user.Role,
	}
}
