package routes

import (
	"github.com/gin-gonic/gin"

	handler "github.com/Aur4ik/AlaRent/internal/handlers"
	"github.com/Aur4ik/AlaRent/internal/middleware"
)

func SetupRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}

	// Bug fix: /me route existed in the handler but was never registered.
	// Added here under the auth middleware so it's actually reachable.
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handler.Me)
	}
}
