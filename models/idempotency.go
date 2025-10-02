package models

type IdempotencyRecord struct {
	Status     string `json:"status"` // "in-progress" or "completed"
	StatusCode int    `json:"status_code,omitempty"`
	Body       string `json:"body,omitempty"`
}
