package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/preraksudan/distributed-job-scheduler/internal/db"
	"github.com/preraksudan/distributed-job-scheduler/internal/models"
	"github.com/preraksudan/distributed-job-scheduler/internal/queue"
)

func main() {

	db.ConnectPostgres()
	queue.ConnectRedis()

	log.Println("Worker service started")

	for {

		processJob()
	}
}

func processJob() {

	result, err := queue.Client.BRPop(
		context.Background(),
		0,
		"job_queue",
	).Result()

	if err != nil {
		log.Println("Redis pop error:", err)
		time.Sleep(2 * time.Second)
		return
	}

	if len(result) < 2 {
		return
	}

	jobJSON := result[1]

	var job models.Job

	err = json.Unmarshal([]byte(jobJSON), &job)

	if err != nil {
		log.Println("Failed to unmarshal job:", err)
		return
	}

	log.Printf("Executing job: %s", job.Name)

	executeJob(job)
}

func executeJob(job models.Job) {

	startedAt := time.Now()

	executionID := uuid.New().String()

	payloadBytes, _ := json.Marshal(job.Payload)

	req, err := http.NewRequest(
		job.Method,
		job.TargetURL,
		bytes.NewBuffer(payloadBytes),
	)

	if err != nil {

		saveExecution(
			executionID,
			job.ID,
			"FAILED",
			startedAt,
			0,
			err.Error(),
		)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Duration(job.TimeoutSeconds) * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {

		saveExecution(
			executionID,
			job.ID,
			"FAILED",
			startedAt,
			0,
			err.Error(),
		)

		return
	}

	defer resp.Body.Close()

	saveExecution(
		executionID,
		job.ID,
		"SUCCESS",
		startedAt,
		resp.StatusCode,
		"",
	)

	log.Printf(
		"Job completed: %s status=%d",
		job.Name,
		resp.StatusCode,
	)
}

func saveExecution(
	executionID string,
	jobID string,
	status string,
	startedAt time.Time,
	responseCode int,
	errorMessage string,
) {

	query := `
	INSERT INTO job_executions (
		id,
		job_id,
		status,
		started_at,
		finished_at,
		response_code,
		error_message
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := db.Conn.Exec(
		context.Background(),
		query,
		executionID,
		jobID,
		status,
		startedAt,
		time.Now(),
		responseCode,
		errorMessage,
	)

	if err != nil {
		log.Println("Failed to save execution:", err)
	}
}
