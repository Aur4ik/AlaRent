package dto

import "github.com/Aur4ik/AlaRent/internal/models"

type RegisterInput struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Group     string `json:"group"`
}


type UserResponse struct {
    ID        uint   `json:"id"`
    Email     string `json:"email"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Group     string `json:"group"`
    Role      string `json:"role"`
    Avatar    string `json:"avatar"`
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Group:     user.Group,
		Role:      user.Role,
		Avatar:    user.Avatar,
	}
}