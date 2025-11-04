package handlers

import (
	"net/http"

	"vigovia-task/models"
	"vigovia-task/services"

	"github.com/gin-gonic/gin"
)

// ItineraryHandler handles HTTP requests for itineraries
type ItineraryHandler struct {
	service    *services.ItineraryService
	pdfService *services.PDFService
}

// NewItineraryHandler creates a new instance of ItineraryHandler
func NewItineraryHandler(service *services.ItineraryService, pdfService *services.PDFService) *ItineraryHandler {
	return &ItineraryHandler{
		service:    service,
		pdfService: pdfService,
	}
}

// CreateItinerary handles POST /itineraries
func (h *ItineraryHandler) CreateItinerary(c *gin.Context) {
	var req models.CreateItineraryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Set the user ID from the authenticated user
	req.UserID = userID.(string)

	itinerary, err := h.service.CreateItinerary(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, itinerary)
}

// GetItinerary handles GET /itineraries/:id
func (h *ItineraryHandler) GetItinerary(c *gin.Context) {
	id := c.Param("id")

	itinerary, err := h.service.GetItinerary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// ListItineraries handles GET /itineraries
func (h *ItineraryHandler) ListItineraries(c *gin.Context) {
	itineraries := h.service.ListItineraries()
	c.JSON(http.StatusOK, gin.H{"itineraries": itineraries})
}

// UpdateItinerary handles PUT /itineraries/:id
func (h *ItineraryHandler) UpdateItinerary(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateItineraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get authenticated user ID from context if available
	userID, exists := c.Get("userID")
	if exists {
		// Override the user_id from request with authenticated user's ID
		req.UserID = userID.(string)
	}

	itinerary, err := h.service.UpdateItinerary(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// DeleteItinerary handles DELETE /itineraries/:id
func (h *ItineraryHandler) DeleteItinerary(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteItinerary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itinerary deleted successfully"})
}

// AddActivity handles POST /itineraries/:id/activities
func (h *ItineraryHandler) AddActivity(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		DayNumber int             `json:"day_number" binding:"required"`
		Activity  models.Activity `json:"activity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itinerary, err := h.service.AddActivity(id, req.DayNumber, &req.Activity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, itinerary)
}

// ExportPDF handles GET /itineraries/:id/export-pdf
func (h *ItineraryHandler) ExportPDF(c *gin.Context) {
	id := c.Param("id")

	itinerary, err := h.service.GetItinerary(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := h.pdfService.GeneratePDF(itinerary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=itinerary.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
