package dtos

// ConfigurationGroupConfigurationDto represents a configuration reference in a group
// swagger:model ConfigurationGroupConfigurationDto
type ConfigurationGroupConfigurationDto struct {
	// The ID of the configuration to include in the group
	// in: string
	// example: config-123
	Id     string            `json:"id"`
	// Labels to apply to this configuration in the group context
	// in: object
	// example: {"environment": "production", "region": "us-east-1"}
	Labels map[string]string `json:"labels"`
}

// CreateConfigurationDto represents the request body for creating a new configuration
// swagger:model CreateConfigurationDto
type CreateConfigurationDto struct {
	// The name of the configuration
	// in: string
	// example: Database Configuration
	Name       string            `json:"name"`
	// The version of the configuration
	// in: string
	// example: v1.0.0
	Version    string            `json:"version"`
	// Key-value pairs representing configuration parameters
	// in: object
	// example: {"host": "localhost", "port": "5432"}
	Parameters map[string]string `json:"parameters"`
}
