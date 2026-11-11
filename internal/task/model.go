package task

import (
	"encoding/json"
	"time"
)

type Task struct {
    ID int64 `json:"id"`
	TaskType string `json:"task_type"`
	Payload json.RawMessage `json:"payload"`
	Status string `json:"status"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Attempts int `json:"attempts"`
	MaxRetries int `json:"max_retries"`
}