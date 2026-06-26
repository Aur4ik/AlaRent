package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Aur4ik/AlaRent/internal/models"
	"github.com/Aur4ik/AlaRent/internal/repository"
)

func AuthMiddleware() gin.HandlerFunc {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		tokenString := c.Query("token")
		if header != "" {
			var found bool
			tokenString, found = strings.CutPrefix(header, "Bearer ")
			if !found {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be Bearer <token>"})
				return
			}
		}
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		// JWT numbers are float64 when decoded from JSON
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
			return
		}
		userID := uint(userIDFloat)

		// Verify the user still exists in the database
		user, err := repository.GetUserByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		slog.Debug("authenticated request", "user_id", userID, "role", user.Role)

		c.Set("user_id", int(userID))
		c.Set("role", user.Role)
		c.Next()
	}
}

func LandlordMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		roleString, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid user role",
			})
			return
		}

		if roleString != models.RoleLandlord {
			slog.Debug("landlord access denied", "role", roleString)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "only landlords can perform this action",
			})
			return
		}
		c.Next()
	}
}

func RequireLandlord() gin.HandlerFunc {
	return LandlordMiddleware()
}
