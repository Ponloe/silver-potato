package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	InitDB()

	// Create Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Movie endpoints
	router.GET("/movies", GetMovies)
	router.GET("/movies/:id", GetMovie)
	router.POST("/movies/:id/decrease", DecreaseSeats)

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
