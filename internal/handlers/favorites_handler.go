package handler

import (
	"net/http"

	"github.com/Aur4ik/AlaRent/internal/service"
	"github.com/gin-gonic/gin"
)

func AddFavorite(c *gin.Context) {
	apartmentID, ok := parseApartmentID(c)
	if !ok {
		return
	}

	if err := service.AddFavorite(uint(c.GetInt("user_id")), apartmentID); err != nil {
		writeApartmentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "apartment added to favorites"})
}

func RemoveFavorite(c *gin.Context) {
	apartmentID, ok := parseApartmentID(c)
	if !ok {
		return
	}

	if err := service.RemoveFavorite(uint(c.GetInt("user_id")), apartmentID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "apartment removed from favorites"})
}

func GetFavorites(c *gin.Context) {
	favorites, err := service.GetUserFavorites(uint(c.GetInt("user_id")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, favorites)
}
