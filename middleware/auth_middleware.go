package middleware

import (
	"net/http"
	"strings"
	"vigovia-task/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is middleware for authenticating requests
func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		userID, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set userID in context for use in handlers
		c.Set("userID", userID)
		c.Next()
	}
}

// OptionalAuthMiddleware is middleware that extracts userID if token is present but doesn't require it
func OptionalAuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token != "" {
			userID, err := authService.ValidateToken(token)
			if err == nil {
				c.Set("userID", userID)
			}
		}
		c.Next()
	}
}

// extractToken extracts the token from the Authorization header or cookie
func extractToken(c *gin.Context) string {
	// First, try to get token from Authorization header
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	
	// If not in header, try to get from cookie
	token, err := c.Cookie("auth_token")
	if err == nil && token != "" {
		return token
	}
	
	return ""
}
