package models

// Configuration represents a single configuration item
// swagger:model Configuration
type Configuration struct {
	// The unique identifier for the configuration
	// in: string
	// example: config-123
	Id         string            `json:"id"`
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

// ConfigurationGroup represents a group of related configurations
// swagger:model ConfigurationGroup
type ConfigurationGroup struct {
	// The unique identifier for the configuration group
	// in: string
	// example: group-456
	Id             string                  `json:"id"`
	// The name of the configuration group
	// in: string
	// example: Database Configurations
	Name           string                  `json:"name"`
	// The version of the configuration group
	// in: string
	// example: v2.1.0
	Version        string                  `json:"version"`
	// Array of labeled configurations in this group
	// in: array
	Configurations []*LabeledConfiguration `json:"configurations"`
}

// LabeledConfiguration represents a configuration with associated labels
// swagger:model LabeledConfiguration
type LabeledConfiguration struct {
	// The unique identifier for the labeled configuration
	// in: string
	// example: labeled-config-789
	Id            string            `json:"id"`
	// The configuration object
	// in: object
	Configuration *Configuration    `json:"configuration"`
	// Key-value pairs representing labels for this configuration
	// in: object
	// example: {"environment": "production", "region": "us-east-1"}
	Labels        map[string]string `json:"labels"`
}
