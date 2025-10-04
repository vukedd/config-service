package models

// ErrorResponse represents error response with status and message
// swagger:model ErrorResponse
type ErrorResponse struct {
	// HTTP status code
	// in: integer
	// example: 404
	Status int `json:"status"`
	// Error message
	// in: string
	// example: Configuration not found
	Message string `json:"message"`
}

// NoContentResponse represents empty response for successful operations with no content
// swagger:model NoContentResponse
type NoContentResponse struct {
	// This struct is intentionally empty for 204 No Content responses
}

