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

	// CORS middleware for Next.js frontend
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// Movie CRUD endpoints
	router.GET("/movies", GetMovies)
	router.POST("/movies", CreateMovie)
	router.GET("/movies/:id", GetMovie)
	router.PUT("/movies/:id", UpdateMovie)
	router.DELETE("/movies/:id", DeleteMovie)
	router.POST("/movies/:id/decrease", DecreaseSeats)

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
