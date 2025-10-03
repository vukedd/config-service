package handlers

import "github.com/vukedd/config-service/dtos"

// CreateConfigurationRequest defines the request body for creating a configuration
// swagger:parameters createConfiguration
type CreateConfigurationRequest struct {
	// Configuration data to create
	// in: body
	// required: true
	Body dtos.CreateConfigurationDto `json:"body"`
}

// CreateConfigurationGroupRequest defines the request body for creating a configuration group
// swagger:parameters createConfigurationGroup
type CreateConfigurationGroupRequest struct {
	// Configuration group data to create
	// in: body
	// required: true
	Body dtos.ConfigurationGroupDto `json:"body"`
}

// UpdateConfigurationGroupRequest defines the request body for updating a configuration group
// swagger:parameters updateConfigurationGroup
type UpdateConfigurationGroupRequest struct {
	// The ID of the configuration group
	// in: path
	// required: true
	ID string `json:"id"`
	// Configuration group data to update
	// in: body
	// required: true
	Body dtos.ConfigurationGroupDto `json:"body"`
}