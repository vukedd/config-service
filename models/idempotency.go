package models

const (
	StatusInProgress = "in-progress"
	StatusCompleted  = "completed"
)

type IdempotencyRecord struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	Body       string `json:"body,omitempty"`
}
