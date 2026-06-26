package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Aur4ik/AlaRent/internal/handlers"
	"github.com/Aur4ik/AlaRent/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "AlaRent API"})
	})
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

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

	registerApartmentRoutes(r.Group("/apartaments"))
	registerApartmentRoutes(r.Group("/apartments"))
}

func registerApartmentRoutes(apartments *gin.RouterGroup) {
	apartments.GET("", handler.GetAllApartments)
	apartments.GET("/:id", handler.GetApartmentByID)
	apartments.POST("", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.CreateApartament)
	apartments.PATCH("/:id", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.UpdateApartment)
	apartments.DELETE("/:id", middleware.AuthMiddleware(), middleware.LandlordMiddleware(), handler.DeleteApartment)
	apartments.POST("/:id/favorite", middleware.AuthMiddleware(), handler.AddFavorite)
	apartments.DELETE("/:id/favorite", middleware.AuthMiddleware(), handler.RemoveFavorite)
	apartments.POST("/:id/conversation", middleware.AuthMiddleware(), handler.StartConversation)
}
