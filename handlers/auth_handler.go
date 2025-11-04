package handlers

import (
	"net/http"
	"strings"
	"vigovia-task/models"
	"vigovia-task/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Signup handles user registration
func (ah *AuthHandler) Signup(c *gin.Context) {
	var req models.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := ah.authService.Signup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set token in cookie for automatic authentication
	c.SetCookie(
		"auth_token",           // name
		authResponse.Token,     // value
		3600*24*7,             // maxAge (7 days in seconds)
		"/",                   // path
		"",                    // domain
		false,                 // secure (set to true in production with HTTPS)
		true,                  // httpOnly
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    authResponse.User.ToUserResponse(),
		"token":   authResponse.Token,
	})
}

// Login handles user authentication
func (ah *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authResponse, err := ah.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set token in cookie for automatic authentication
	c.SetCookie(
		"auth_token",           // name
		authResponse.Token,     // value
		3600*24*7,             // maxAge (7 days in seconds)
		"/",                   // path
		"",                    // domain
		false,                 // secure (set to true in production with HTTPS)
		true,                  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    authResponse.User.ToUserResponse(),
		"token":   authResponse.Token,
	})
}

// Logout handles user logout
func (ah *AuthHandler) Logout(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	if err := ah.authService.Logout(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	// Clear the auth cookie
	c.SetCookie(
		"auth_token",  // name
		"",           // value
		-1,           // maxAge (negative value deletes cookie)
		"/",          // path
		"",           // domain
		false,        // secure
		true,         // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// GetProfile returns the current user's profile
func (ah *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := ah.authService.GetUserByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToUserResponse(),
	})
}

// extractToken extracts the token from the Authorization header
func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
