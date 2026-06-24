package dto

import "github.com/Aur4ik/AlaRent/internal/models"

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"    binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"     binding:"omitempty,oneof=tenant landlord"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UpdateProfileRequest struct {
	Name      *string `json:"name"`
	Phone     *string `json:"phone"`
	AvatarURL *string `json:"avatar_url"`
	Bio       *string `json:"bio"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

type CreateApartmentRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	Type         string   `json:"type" binding:"omitempty,oneof=apartment room studio house"`
	Price        int      `json:"price" binding:"required,gt=0"`
	District     string   `json:"district" binding:"required"`
	Address      string   `json:"address" binding:"required"`
	Rooms        int      `json:"rooms" binding:"required,gt=0"`
	Floor        int      `json:"floor" binding:"required,gte=0"`
	HasFurniture bool     `json:"has_furniture"`
	HasWifi      bool     `json:"has_wifi"`
	HasWasher    bool     `json:"has_washer"`
	PhotoURLs    []string `json:"photo_urls" binding:"omitempty,max=10,dive,url"`
}

type UpdateApartmentRequest struct {
	Title        *string   `json:"title"`
	Description  *string   `json:"description"`
	Type         *string   `json:"type" binding:"omitempty,oneof=apartment room studio house"`
	Price        *int      `json:"price" binding:"omitempty,gt=0"`
	District     *string   `json:"district"`
	Address      *string   `json:"address"`
	Rooms        *int      `json:"rooms" binding:"omitempty,gt=0"`
	Floor        *int      `json:"floor" binding:"omitempty,gte=0"`
	HasFurniture *bool     `json:"has_furniture"`
	HasWifi      *bool     `json:"has_wifi"`
	HasWasher    *bool     `json:"has_washer"`
	PhotoURLs    *[]string `json:"photo_urls" binding:"omitempty,max=10,dive,url"`
}

type ApartmentFilter struct {
	Query        string
	District     string
	Type         string
	MinPrice     int
	MaxPrice     int
	Rooms        int
	HasFurniture *bool
	HasWifi      *bool
	HasWasher    *bool
	Sort         string
}

type SendMessageRequest struct {
	Text string `json:"text" binding:"required"`
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
		Bio:       user.Bio,
	}
}
