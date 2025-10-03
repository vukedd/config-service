package models

const (
	StatusInProgress = "in-progress"
	StatusCompleted  = "completed"
)

// IdempotencyRecord represents the status and result of an idempotent operation
// swagger:model IdempotencyRecord
type IdempotencyRecord struct {
	// The current status of the operation
	// in: string
	// enum: in-progress,completed
	// example: completed
	Status string `json:"status"` // "in-progress" or "completed"
	// The HTTP status code of the completed operation
	// in: integer
	// example: 200
	StatusCode int `json:"status_code,omitempty"`
	// The response body of the completed operation
	// in: string
	// example: {"id": "config-123", "created": true}
	Body string `json:"body,omitempty"`
}
