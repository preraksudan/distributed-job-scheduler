package models

import "time"

type Job struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	CronExpression    string                 `json:"cron_expression"`
	TargetURL         string                 `json:"target_url"`
	Method            string                 `json:"method"`
	Payload           map[string]interface{} `json:"payload"`
	Retries           int                    `json:"retries"`
	TimeoutSeconds    int                    `json:"timeout_seconds"`
	Enabled           bool                   `json:"enabled"`
	CreatedAt         time.Time              `json:"created_at"`
	LastExecutionTime *time.Time             `json:"last_execution_time"`
}
