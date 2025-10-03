package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/dtos"
	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
)

type ConfigurationGroupHandler struct {
	gr *repositories.ConfigurationGroupRepository
	cr *repositories.ConfigurationRepository
}

func NewConfigurationGroupHandler(gr *repositories.ConfigurationGroupRepository, cr *repositories.ConfigurationRepository) *ConfigurationGroupHandler {
	return &ConfigurationGroupHandler{
		gr: gr,
		cr: cr,
	}
}

// Helper function to send error response
func (h ConfigurationGroupHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{Status: statusCode, Message: message}
	json.NewEncoder(w).Encode(errorResponse)
}

// FindAll retrieves all configuration groups
// swagger:route GET /configurationGroups configurationGroups getAllConfigurationGroups
//
// # Get all configuration groups
//
// This endpoint retrieves all configuration groups in the system.
//
// Responses:
//
//	200: body:[]ConfigurationGroup
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	configurationGroups, err := h.gr.FindAll()
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(configurationGroups)

	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// FindById retrieves a configuration group by ID
// swagger:route GET /configurationGroups/{id} configurationGroups getConfigurationGroupById
//
// # Get configuration group by ID
//
// This endpoint retrieves a specific configuration group by its ID.
//
// Parameters:
//   - name: id
//     in: path
//     type: string
//     required: true
//     description: The ID of the configuration group
//
// Responses:
//
//	200: body:ConfigurationGroup
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]

	configurationGroup, err := h.gr.FindById(configurationGroupId)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(configurationGroup)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Delete removes a configuration group by ID
// swagger:route DELETE /configurationGroups/{id} configurationGroups deleteConfigurationGroupById
//
// # Delete configuration group by ID
//
// This endpoint deletes a specific configuration group by its ID.
//
// Parameters:
//   - name: id
//     in: path
//     type: string
//     required: true
//     description: The ID of the configuration group
//
// Responses:
//
//	204: body:NoContentResponse
//	404: body:ErrorResponse
func (h ConfigurationGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]
	err := h.gr.DeleteById(configurationGroupId)

	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(models.NoContentResponse{})
}

// Create creates a new configuration group
// swagger:route POST /configurationGroups configurationGroups createConfigurationGroup
//
// # Create a new configuration group
//
// This endpoint creates a new configuration group with the provided data.
//
// Responses:
//
//	200: body:ConfigurationGroup
//	400: body:ErrorResponse
//	404: body:ErrorResponse
//	409: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var configurationGroupRequest dtos.ConfigurationGroupDto
	_ = json.NewDecoder(r.Body).Decode(&configurationGroupRequest)

	if len(configurationGroupRequest.ConfigurationList) < 1 {
		h.sendErrorResponse(w, http.StatusBadRequest, "you must define at least one configuration")
		return
	}

	_, err := h.gr.FindByNameAndVersion(configurationGroupRequest.Name, configurationGroupRequest.Version)
	if err == nil {
		h.sendErrorResponse(w, http.StatusConflict, "configuration group already exists")
		return
	}

	if err != repositories.ErrConfigurationGroupNotFound {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// I thought it was a good idea to leave the transformation from dto the model in the handler since we are going
	// to check if selected configurations for this configuration group exist by fetching data from the repository,
	// and by leaving it as it is, I am avoiding giving mapper classes data access :D SAME GOES FOR THE UPDATE METHOD
	configurationGroupConfigurationList := []*models.LabeledConfiguration{}
	for _, configurationItem := range configurationGroupRequest.ConfigurationList {
		configuration, err := h.cr.FindById(configurationItem.Id)

		if err == repositories.ErrConfigurationNotFound {
			h.sendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("configuration with the id %s does not exist", configurationItem.Id))
			return
		}

		if err != nil {
			h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		configurationGroupConfigurationList = append(configurationGroupConfigurationList, &models.LabeledConfiguration{Id: "", Configuration: configuration, Labels: configurationItem.Labels})
	}

	newConfigurationGroup := models.ConfigurationGroup{Id: "", Name: configurationGroupRequest.Name, Version: configurationGroupRequest.Version, Configurations: configurationGroupConfigurationList}
	configGroup, err := h.gr.Create(&newConfigurationGroup)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(configGroup)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Update updates an existing configuration group
// swagger:route PUT /configurationGroups/{id} configurationGroups updateConfigurationGroup
//
// # Update configuration group by ID
//
// This endpoint updates an existing configuration group with the provided data.
//
// Responses:
//
//	200: body:ConfigurationGroup
//	400: body:ErrorResponse
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]

	var configGroupData dtos.ConfigurationGroupDto
	json.NewDecoder(r.Body).Decode(&configGroupData)

	if len(configGroupData.ConfigurationList) < 1 {
		h.sendErrorResponse(w, http.StatusBadRequest, "you must define at least one configuration")
		return
	}

	labeledConfiguration := []*models.LabeledConfiguration{}
	for _, ci := range configGroupData.ConfigurationList {
		c, err := h.cr.FindById(configurationGroupId)
		if err == repositories.ErrConfigurationGroupNotFound {
			h.sendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("configuration with the id %s does not exist", configurationGroupId))
			return
		}

		if err != nil {
			h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		labeledConfiguration = append(labeledConfiguration, &models.LabeledConfiguration{Id: "", Configuration: c, Labels: ci.Labels})
	}

	updateConfigurationGroup := models.ConfigurationGroup{Id: configurationGroupId, Name: configGroupData.Name, Version: configGroupData.Version, Configurations: labeledConfiguration}
	err := h.gr.Update(configurationGroupId, &updateConfigurationGroup)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// FindByIdToDto retrieves a configuration group by ID as DTO
// swagger:route GET /configurationGroups/dto/{id} configurationGroups getConfigurationGroupByIdToDto
//
// # Get configuration group by ID as DTO
//
// This endpoint retrieves a specific configuration group by its ID and returns it as a DTO.
//
// Parameters:
//   - name: id
//     in: path
//     type: string
//     required: true
//     description: The ID of the configuration group
//
// Responses:
//
//	200: body:ConfigurationGroupDto
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) FindByIdToDto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]

	configurationGroup, err := h.gr.FindById(configurationGroupId)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	configurationsTransformedToDto := []*dtos.ConfigurationGroupConfigurationDto{}
	for _, labeledConfiguration := range configurationGroup.Configurations {
		configurationsTransformedToDto = append(configurationsTransformedToDto, &dtos.ConfigurationGroupConfigurationDto{Id: labeledConfiguration.Configuration.Id, Labels: labeledConfiguration.Labels})
	}

	configurationGroupDto := dtos.ConfigurationGroupDto{Name: configurationGroup.Name, Version: configurationGroup.Version, ConfigurationList: configurationsTransformedToDto}
	err = json.NewEncoder(w).Encode(configurationGroupDto)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// FindByNameAndVersion retrieves a configuration group by name and version
// swagger:route GET /configurationGroups/{name}/{version} configurationGroups getConfigurationGroupByNameAndVersion
//
// # Get configuration group by name and version
//
// This endpoint retrieves a specific configuration group by its name and version.
//
// Parameters:
//   - name: name
//     in: path
//     type: string
//     required: true
//     description: The name of the configuration group
//   - name: version
//     in: path
//     type: string
//     required: true
//     description: The version of the configuration group
//
// Responses:
//
//	200: body:ConfigurationGroup
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) FindByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupName := params["name"]
	configurationGroupVersion := params["version"]

	configurationGroup, err := h.gr.FindByNameAndVersion(configurationGroupName, configurationGroupVersion)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(configurationGroup)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// FindByNameAndVersionToDto retrieves a configuration group by name and version as DTO
// swagger:route GET /configurationGroups/dto/{name}/{version} configurationGroups getConfigurationGroupByNameAndVersionToDto
//
// # Get configuration group by name and version as DTO
//
// This endpoint retrieves a specific configuration group by its name and version and returns it as a DTO.
//
// Parameters:
//   - name: name
//     in: path
//     type: string
//     required: true
//     description: The name of the configuration group
//   - name: version
//     in: path
//     type: string
//     required: true
//     description: The version of the configuration group
//
// Responses:
//
//	200: body:ConfigurationGroupDto
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationGroupHandler) FindByNameAndVersionToDto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupName := params["name"]
	configurationGroupVersion := params["version"]

	configurationGroup, err := h.gr.FindByNameAndVersion(configurationGroupName, configurationGroupVersion)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	configurationsTransformedToDto := []*dtos.ConfigurationGroupConfigurationDto{}
	for _, labeledConfiguration := range configurationGroup.Configurations {
		configurationsTransformedToDto = append(configurationsTransformedToDto, &dtos.ConfigurationGroupConfigurationDto{Id: labeledConfiguration.Configuration.Id, Labels: labeledConfiguration.Labels})
	}

	configurationGroupDto := dtos.ConfigurationGroupDto{Name: configurationGroup.Name, Version: configurationGroup.Version, ConfigurationList: configurationsTransformedToDto}
	err = json.NewEncoder(w).Encode(configurationGroupDto)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// DeleteByNameAndVersion removes a configuration group by name and version
// swagger:route DELETE /configurationGroups/{name}/{version} configurationGroups deleteConfigurationGroupByNameAndVersion
//
// # Delete configuration group by name and version
//
// This endpoint deletes a specific configuration group by its name and version.
//
// Parameters:
//   - name: name
//     in: path
//     type: string
//     required: true
//     description: The name of the configuration group
//   - name: version
//     in: path
//     type: string
//     required: true
//     description: The version of the configuration group
//
// Responses:
//
//	204: body:NoContentResponse
//	404: body:ErrorResponse
func (h ConfigurationGroupHandler) DeleteByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupName := params["name"]
	configurationGroupVersion := params["version"]

	err := h.gr.DeleteByNameAndVersion(configurationGroupName, configurationGroupVersion)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(models.NoContentResponse{})
}
