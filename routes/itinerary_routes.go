package routes

import (
	"vigovia-task/handlers"
	"vigovia-task/middleware"
	"vigovia-task/services"
	"vigovia-task/storage"

	"github.com/gin-gonic/gin"
)

// RegisterItineraryRoutes registers all itinerary-related routes
func RegisterItineraryRoutes(router *gin.Engine) {
	// Initialize storage, services, and handlers
	store := storage.NewMemoryStore()
	
	// Auth services and handlers
	authService := services.NewAuthService(store)
	authHandler := handlers.NewAuthHandler(authService)
	
	// Itinerary services and handlers
	itineraryService := services.NewItineraryService(store)
	pdfService := services.NewPDFService()
	itineraryHandler := handlers.NewItineraryHandler(itineraryService, pdfService)

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/profile", middleware.AuthMiddleware(authService), authHandler.GetProfile)
		}

		// Itinerary routes (protected)
		itineraries := api.Group("/itineraries")
		itineraries.Use(middleware.AuthMiddleware(authService))
		{
			itineraries.POST("", itineraryHandler.CreateItinerary)
			itineraries.GET("", itineraryHandler.ListItineraries)
			itineraries.GET("/:id", itineraryHandler.GetItinerary)
			itineraries.PUT("/:id", itineraryHandler.UpdateItinerary)
			itineraries.DELETE("/:id", itineraryHandler.DeleteItinerary)
			itineraries.POST("/:id/activities", itineraryHandler.AddActivity)
			itineraries.GET("/:id/export-pdf", itineraryHandler.ExportPDF)
		}
	}

	// Health check and welcome routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Itinerary API is running successfully!",
			"version": "1.0.0",
			
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "Itinerary Builder API",
			"version": "1.0.0",
		})
	})
}
