package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMovies returns all movies
func GetMovies(c *gin.Context) {
	var movies []Movie
	if err := DB.Find(&movies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch movies",
		})
		return
	}

	c.JSON(http.StatusOK, movies)
}

// GetMovie returns a single movie by ID
func GetMovie(c *gin.Context) {
	id := c.Param("id")

	var movie Movie
	if err := DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found",
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// CreateMovie creates a new movie
func CreateMovie(c *gin.Context) {
	var movie Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if movie.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title is required",
		})
		return
	}

	if movie.AvailableSeats < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Available seats must be non-negative",
		})
		return
	}

	if err := DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create movie",
		})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// UpdateMovie updates an existing movie
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var movie Movie
	if err := DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found",
		})
		return
	}

	var updateData struct {
		Title          *string `json:"title"`
		AvailableSeats *int    `json:"available_seats"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if updateData.Title != nil {
		movie.Title = *updateData.Title
	}
	if updateData.AvailableSeats != nil {
		movie.AvailableSeats = *updateData.AvailableSeats
	}

	if err := DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update movie",
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// DeleteMovie deletes a movie
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	if err := DB.Delete(&Movie{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete movie",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully",
	})
}

// DecreaseSeats decreases available seats for a movie
func DecreaseSeats(c *gin.Context) {
	id := c.Param("id")

	var requestBody struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body. 'quantity' must be a positive integer",
		})
		return
	}

	var movie Movie
	if err := DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found",
		})
		return
	}

	// Check if enough seats are available
	if movie.AvailableSeats < requestBody.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Not enough seats available",
			"available_seats": movie.AvailableSeats,
			"requested":       requestBody.Quantity,
		})
		return
	}

	// Decrease seats
	movie.AvailableSeats -= requestBody.Quantity
	if err := DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update movie",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Seats decreased successfully",
		"movie_id":        movie.ID,
		"available_seats": movie.AvailableSeats,
		"decreased_by":    requestBody.Quantity,
	})
}
