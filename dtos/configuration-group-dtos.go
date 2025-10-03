package dtos

// ConfigurationGroupDto represents the request/response body for configuration groups
// swagger:model ConfigurationGroupDto
type ConfigurationGroupDto struct {
	// The name of the configuration group
	// in: string
	// example: Database Configurations
	Name              string                                `json:"name"`
	// The version of the configuration group
	// in: string
	// example: v2.1.0
	Version           string                                `json:"version"`
	// List of configurations to include in this group
	// in: array
	// minItems: 1
	ConfigurationList []*ConfigurationGroupConfigurationDto `json:"configuration_list"`
}
