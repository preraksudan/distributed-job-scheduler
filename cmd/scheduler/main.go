package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/preraksudan/distributed-job-scheduler/internal/db"
	"github.com/preraksudan/distributed-job-scheduler/internal/models"
	"github.com/preraksudan/distributed-job-scheduler/internal/queue"

	"github.com/robfig/cron/v3"
)

func main() {

	db.ConnectPostgres()
	queue.ConnectRedis()

	log.Println("Scheduler service started")

	for {

		checkAndQueueJobs()

		time.Sleep(10 * time.Second)
	}
}

func checkAndQueueJobs() {

	query := `
	SELECT
		id,
		name,
		cron_expression,
		target_url,
		method,
		payload,
		retries,
		timeout_seconds,
		enabled,
		created_at
	FROM jobs
	WHERE enabled = true
	`

	rows, err := db.Conn.Query(context.Background(), query)

	if err != nil {
		log.Println("Error querying jobs:", err)
		return
	}

	defer rows.Close()

	for rows.Next() {

		var job models.Job

		err := rows.Scan(
			&job.ID,
			&job.Name,
			&job.CronExpression,
			&job.TargetURL,
			&job.Method,
			&job.Payload,
			&job.Retries,
			&job.TimeoutSeconds,
			&job.Enabled,
			&job.CreatedAt,
		)

		if err != nil {
			log.Println("Error scanning job:", err)
			continue
		}

		parser := cron.NewParser(
			cron.Minute |
				cron.Hour |
				cron.Dom |
				cron.Month |
				cron.Dow,
		)

		schedule, err := parser.Parse(job.CronExpression)

		if err != nil {
			log.Println("Invalid cron:", err)
			continue
		}

		now := time.Now()

		nextRun := schedule.Next(now.Add(-1 * time.Minute))

		if nextRun.Before(now) || nextRun.Equal(now) {

			jobJSON, _ := json.Marshal(job)

			err := queue.Client.LPush(
				context.Background(),
				"job_queue",
				jobJSON,
			).Err()

			if err != nil {
				log.Println("Failed to enqueue job:", err)
				continue
			}

			log.Printf(
				"Queued job: %s (%s)",
				job.Name,
				job.ID,
			)
		}
	}
}
