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
		auth.POST("/refresh", handler.Refresh)
		auth.POST("/logout", handler.Logout)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handler.Me)
		protected.PATCH("/me", handler.UpdateMe)
		protected.GET("/me/favorites", handler.GetFavorites)
		protected.GET("/conversations", handler.GetConversations)
		protected.GET("/conversations/:id/messages", handler.GetMessages)
		protected.POST("/conversations/:id/messages", handler.SendMessage)
		protected.GET("/ws/conversations/:id", handler.ConversationWebSocket)
	}

	apartaments := r.Group("/apartaments")
	{
		apartaments.GET("", handler.GetAllApartments)
		apartaments.GET("/:id", handler.GetApartmentByID)
		apartaments.POST("", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.CreateApartament)
		apartaments.PATCH("/:id", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.UpdateApartment)
		apartaments.DELETE("/:id", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.DeleteApartment)
		apartaments.POST("/:id/favorite", middleware.AuthMiddleware(), handler.AddFavorite)
		apartaments.DELETE("/:id/favorite", middleware.AuthMiddleware(), handler.RemoveFavorite)
		apartaments.POST("/:id/conversation", middleware.AuthMiddleware(), handler.StartConversation)
	}
}
