package handler

import (
	"net/http"

	"github.com/Aur4ik/AlaRent/internal/dto"
	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/service"
	"github.com/gin-gonic/gin"
)

func CreateApartament(c *gin.Context) {
	var req dto.CreateApartmentRequest
	userID := c.GetInt("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apartment := models.Apartment{
		OwnerID:      uint(userID),
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		District:     req.District,
		Address:      req.Address,
		Rooms:        req.Rooms,
		Floor:        req.Floor,
		HasFurniture: req.HasFurniture,
		HasWifi:      req.HasWifi,
	}

	if err := service.CreateAppartament(&apartment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, apartment)
}
func GetAllApartments(c *gin.Context) {

	apartments, err := service.GetAllApartments()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, apartments)
}