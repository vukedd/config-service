package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/dtos"
	"github.com/vukedd/config-service/mappers"
	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
)

type ConfigurationHandler struct {
	repository *repositories.ConfigurationRepository
}

func NewConfigurationHandler(repository *repositories.ConfigurationRepository) *ConfigurationHandler {
	return &ConfigurationHandler{
		repository: repository,
	}
}

// Helper function to send error response
func (h ConfigurationHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{Status: statusCode, Message: message}
	json.NewEncoder(w).Encode(errorResponse)
}

// FindAll retrieves all configurations
// swagger:route GET /configurations configurations getAllConfigurations
//
// Get all configurations
//
// This endpoint retrieves all configurations in the system.
//
// Responses:
//   200: []Configuration
//   500: ErrorResponse
func (h ConfigurationHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	configurations := h.repository.FindAll()

	err := json.NewEncoder(w).Encode(configurations)

	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// FindById retrieves a configuration by ID
// swagger:route GET /configurations/{id} configurations getConfigurationById
//
// Get configuration by ID
//
// This endpoint retrieves a specific configuration by its ID.
//
// Parameters:
//   + name: ConfigurationByIdParams
//
// Responses:
//   200: Configuration
//   404: ErrorResponse
//   500: ErrorResponse
//   404: ErrorResponse
//   500: ErrorResponse
func (h ConfigurationHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	configuration, err := h.repository.FindById(id)

	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(configuration)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Create creates a new configuration
// swagger:route POST /configurations configurations createConfiguration
//
// Create a new configuration
//
// This endpoint creates a new configuration with the provided data.
//
// Parameters:
//   + name: CreateConfigurationParams
//
// Responses:
//   200: Configuration
//   400: ErrorResponse
//   409: ErrorResponse
//   500: ErrorResponse
func (h ConfigurationHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var configuration dtos.CreateConfigurationDto
	_ = json.NewDecoder(r.Body).Decode(&configuration)

	if len(configuration.Parameters) < 1 {
		h.sendErrorResponse(w, http.StatusBadRequest, "you must specify at least one parameter")
		return
	}

	createdConfig, err := h.repository.Create(*(mappers.ToConfiguration(&configuration)))
	if err != nil {
		h.sendErrorResponse(w, http.StatusConflict, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(createdConfig)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// Delete removes a configuration by ID
// swagger:route DELETE /configurations/{id} configurations deleteConfigurationById
//
// Delete configuration by ID
//
// This endpoint deletes a specific configuration by its ID.
//
// Parameters:
//   + name: ConfigurationByIdParams
//
// Responses:
//   204: NoContentResponse
//   404: ErrorResponse
func (h ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.repository.Delete(id)

	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(models.NoContentResponse{})
}

// DeleteByNameAndVersion removes a configuration by name and version
// swagger:route DELETE /configuration/{name}/{version} configurations deleteConfigurationByNameAndVersion
//
// Delete configuration by name and version
//
// This endpoint deletes a specific configuration by its name and version.
//
// Parameters:
//   + name: ConfigurationByNameAndVersionParams
//
// Responses:
//   204: NoContentResponse
//   404: ErrorResponse
func (h ConfigurationHandler) DeleteByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	name := params["name"]
	version := params["version"]

	err := h.repository.DeleteByNameAndVersion(name, version)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(models.NoContentResponse{})
}

// FindByNameAndVersion retrieves a configuration by name and version
// swagger:route GET /configuration/{name}/{version} configurations getConfigurationByNameAndVersion
//
// Get configuration by name and version
//
// This endpoint retrieves a specific configuration by its name and version.
//
// Parameters:
//   + name: ConfigurationByNameAndVersionParams
//
// Responses:
//   200: Configuration
//   404: ErrorResponse
//   500: ErrorResponse
func (h ConfigurationHandler) FindByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	name := params["name"]
	version := params["version"]

	configuration, err := h.repository.FindByNameAndVersion(name, version)
	if err != nil {
		h.sendErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(configuration)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}
