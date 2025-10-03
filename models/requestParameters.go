package models

// in: path indicates the parameter comes from the URL path (e.g., /config/{id})
// in: body indicates the parameter is in the request body (for POST/PUT operations)

import "github.com/vukedd/config-service/dtos"

// ConfigurationByIdParams represents path parameters for operations on a single configuration by ID
// swagger:parameters getConfigurationById deleteConfigurationById
type ConfigurationByIdParams struct {
	// The ID of the configuration
	// in: path
	// required: true
	// example: 12345
	ID string `json:"id"`
}

// ConfigurationByNameAndVersionParams represents path parameters for operations on a configuration by name and version
// swagger:parameters getConfigurationByNameAndVersion deleteConfigurationByNameAndVersion
type ConfigurationByNameAndVersionParams struct {
	// The name of the configuration
	// in: path
	// required: true
	// example: database-config
	Name string `json:"name"`
	// The version of the configuration
	// in: path
	// required: true
	// example: v1.0.0
	Version string `json:"version"`
}

// ConfigurationGroupByIdParams represents path parameters for operations on a single configuration group by ID
// swagger:parameters getConfigurationGroupById deleteConfigurationGroupById getConfigurationGroupByIdToDto updateConfigurationGroup
type ConfigurationGroupByIdParams struct {
	// The ID of the configuration group
	// in: path
	// required: true
	// example: 67890
	ID string `json:"id"`
}

// ConfigurationGroupByNameAndVersionParams represents path parameters for operations on a configuration group by name and version
// swagger:parameters getConfigurationGroupByNameAndVersion deleteConfigurationGroupByNameAndVersion getConfigurationGroupByNameAndVersionToDto
type ConfigurationGroupByNameAndVersionParams struct {
	// The name of the configuration group
	// in: path
	// required: true
	// example: web-app-config
	Name string `json:"name"`
	// The version of the configuration group
	// in: path
	// required: true
	// example: v2.1.0
	Version string `json:"version"`
}

// CreateConfigurationParams represents the request body for creating a new configuration
// swagger:parameters createConfiguration
type CreateConfigurationParams struct {
	// The configuration data to create
	// in: body
	// required: true
	// schema:
	//   $ref: "#/definitions/CreateConfigurationDto"
	Body dtos.CreateConfigurationDto `json:"body"`
}

// CreateConfigurationGroupParams represents the request body for creating a new configuration group
// swagger:parameters createConfigurationGroup
type CreateConfigurationGroupParams struct {
	// The configuration group data to create
	// in: body
	// required: true
	// schema:
	//   $ref: "#/definitions/ConfigurationGroupDto"
	Body dtos.ConfigurationGroupDto `json:"body"`
}

// UpdateConfigurationGroupParams represents parameters for updating a configuration group
// swagger:parameters updateConfigurationGroup
type UpdateConfigurationGroupParams struct {
	// The updated configuration group data
	// in: body
	// required: true
	// schema:
	//   $ref: "#/definitions/ConfigurationGroupDto"
	Body dtos.ConfigurationGroupDto `json:"body"`
}