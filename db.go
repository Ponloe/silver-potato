package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default values")
	}

	// Get environment variables or use defaults
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "inventory")

	// Create connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Connect to database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	// Auto migrate the schema
	err = DB.AutoMigrate(&Movie{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed successfully")

	// Seed initial data
	SeedMovies()
}

func SeedMovies() {
	var count int64
	DB.Model(&Movie{}).Count(&count)

	if count == 0 {
		log.Println("Seeding initial movies...")
		movies := []Movie{
			{Title: "The Matrix", AvailableSeats: 50},
			{Title: "Inception", AvailableSeats: 40},
			{Title: "Interstellar", AvailableSeats: 45},
			{Title: "The Dark Knight", AvailableSeats: 60},
		}

		for _, movie := range movies {
			if err := DB.Create(&movie).Error; err != nil {
				log.Printf("Failed to seed movie %s: %v", movie.Title, err)
			} else {
				log.Printf("Seeded movie: %s", movie.Title)
			}
		}
		log.Println("Movie seeding completed")
	} else {
		log.Printf("Database already contains %d movies, skipping seed", count)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
