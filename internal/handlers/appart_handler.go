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
		Type:         req.Type,
		Price:        req.Price,
		District:     req.District,
		Address:      req.Address,
		Rooms:        req.Rooms,
		Floor:        req.Floor,
		HasFurniture: req.HasFurniture,
		HasWifi:      req.HasWifi,
		HasWasher:    req.HasWasher,
	}

	if err := service.CreateApartment(&apartment, req.PhotoURLs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, apartment)
}
func GetAllApartments(c *gin.Context) {
	filter, ok := parseApartmentFilter(c)
	if !ok {
		return
	}

	apartments, err := service.GetAllApartments(filter)

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

func parseApartmentFilter(c *gin.Context) (dto.ApartmentFilter, bool) {
	filter := dto.ApartmentFilter{
		Query:    c.Query("q"),
		District: c.Query("district"),
		Type:     c.Query("type"),
		Sort:     c.Query("sort"),
	}

	var ok bool
	if filter.MinPrice, ok = parseOptionalInt(c, "min_price"); !ok {
		return filter, false
	}
	if filter.MaxPrice, ok = parseOptionalInt(c, "max_price"); !ok {
		return filter, false
	}
	if filter.Rooms, ok = parseOptionalInt(c, "rooms"); !ok {
		return filter, false
	}
	if filter.HasFurniture, ok = parseOptionalBool(c, "has_furniture"); !ok {
		return filter, false
	}
	if filter.HasWifi, ok = parseOptionalBool(c, "has_wifi"); !ok {
		return filter, false
	}
	if filter.HasWasher, ok = parseOptionalBool(c, "has_washer"); !ok {
		return filter, false
	}

	return filter, true
}

func parseOptionalInt(c *gin.Context, key string) (int, bool) {
	raw := c.Query(key)
	if raw == "" {
		return 0, true
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return 0, false
	}

	return value, true
}

func parseOptionalBool(c *gin.Context, key string) (*bool, bool) {
	raw := c.Query(key)
	if raw == "" {
		return nil, true
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return nil, false
	}

	return &value, true
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
