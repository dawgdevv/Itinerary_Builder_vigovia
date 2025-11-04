package main

import (
	"log"

	"vigovia-task/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	routes.RegisterItineraryRoutes(router)

	log.Printf("Starting Itinerary Builder API on http://localhost:8080\n")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
