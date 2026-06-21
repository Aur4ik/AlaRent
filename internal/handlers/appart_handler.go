package handler

import (
	"errors"
	"net/http"
	"strconv"

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

func GetApartmentByID(c *gin.Context) {
	apartmentID, ok := parseApartmentID(c)
	if !ok {
		return
	}

	apartment, err := service.GetApartmentByID(apartmentID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "apartment not found",
		})
		return
	}
	c.JSON(http.StatusOK, apartment)
}

func UpdateApartment(c *gin.Context) {
	apartmentID, ok := parseApartmentID(c)
	if !ok {
		return
	}

	var req dto.UpdateApartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apartment, err := service.UpdateApartment(apartmentID, uint(c.GetInt("user_id")), req)
	if err != nil {
		writeApartmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, apartment)
}

func DeleteApartment(c *gin.Context) {
	apartmentID, ok := parseApartmentID(c)
	if !ok {
		return
	}

	if err := service.DeleteApartment(apartmentID, uint(c.GetInt("user_id"))); err != nil {
		writeApartmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "apartment deleted"})
}

func parseApartmentID(c *gin.Context) (uint, bool) {
	apartmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid apartment id",
		})
		return 0, false
	}

	return uint(apartmentID), true
}

func writeApartmentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrApartmentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrApartmentForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
