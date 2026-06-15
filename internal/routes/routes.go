package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Aur4ik/AlaRent/internal/handlers"
	"github.com/Aur4ik/AlaRent/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handler.Me)
	}

	apartaments := r.Group("/apartaments")
	{
		apartaments.GET("", handler.GetAllApartments)
		apartaments.POST("", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.CreateApartament)
	}
}
