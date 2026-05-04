package jobs

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/preraksudan/distributed-job-scheduler/internal/db"
	"github.com/preraksudan/distributed-job-scheduler/internal/models"
	"github.com/robfig/cron/v3"
)

func CreateJob(c *gin.Context) {

	var job models.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate cron expression
	_, err := cron.ParseStandard(job.CronExpression)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid cron expression",
		})
		return
	}

	job.ID = uuid.New().String()

	query := `
	INSERT INTO jobs (
		id,
		name,
		cron_expression,
		target_url,
		method,
		payload,
		retries,
		timeout_seconds,
		enabled
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`

	_, err = db.Conn.Exec(
		context.Background(),
		query,
		job.ID,
		job.Name,
		job.CronExpression,
		job.TargetURL,
		job.Method,
		job.Payload,
		job.Retries,
		job.TimeoutSeconds,
		true,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "job created",
		"job_id":  job.ID,
	})
}
