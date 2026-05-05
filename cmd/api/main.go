package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/preraksudan/distributed-job-scheduler/internal/db"
	"github.com/preraksudan/distributed-job-scheduler/internal/jobs"
)

func main() {
	// 1. Load the .env file from the root directory
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	// 2. db init
	db.ConnectPostgres()

	// 3. Setup Router
	router := gin.Default()

	// Endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API running",
		})
	})

	router.POST("/jobs", jobs.CreateJob)

	// 4. Handle Port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8082" // Default fallback
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
