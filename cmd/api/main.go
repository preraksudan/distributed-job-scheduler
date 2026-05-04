package main

import (
	"github.com/gin-gonic/gin"
	"github.com/preraksudan/distributed-job-scheduler/internal/db"
	"github.com/preraksudan/distributed-job-scheduler/internal/jobs"
)

func main() {

	db.ConnectPostgres()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API running",
		})
	})

	router.POST("/jobs", jobs.CreateJob)

	router.Run(":8080")
}
