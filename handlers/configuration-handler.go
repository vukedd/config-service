package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/dtos"
	"github.com/vukedd/config-service/mappers"
	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ConfigurationHandler struct {
    r      *repositories.ConfigurationRepository
    Tracer trace.Tracer
}

func NewConfigurationHandler(r *repositories.ConfigurationRepository, tracer trace.Tracer) *ConfigurationHandler {
    return &ConfigurationHandler{
        r:      r,
        Tracer: tracer,
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
// # Get all configurations
//
// This endpoint retrieves all configurations in the system.
//
// Responses:
//
//	200: body:[]Configuration
//	500: body:ErrorResponse
func (h ConfigurationHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "ConfigurationHandler.FindAll")
    defer span.End()

    w.Header().Set("Content-Type", "application/json")
    configurations, err := h.r.FindAll(ctx)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
        return
    }

    err = json.NewEncoder(w).Encode(configurations)

    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
    } else {
        span.SetStatus(codes.Ok, "")
    }
}

// FindById retrieves a configuration by ID
// swagger:route GET /configurations/{id} configurations getConfigurationById
//
// # Get configuration by ID
//
// This endpoint retrieves a specific configuration by its ID.
//
// Parameters:
//   - name: id
//     in: path
//     type: string
//     required: true
//     description: The ID of the configuration
//
// Responses:
//
//	200: body:Configuration
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationHandler) FindById(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "ConfigurationHandler.FindById")
    defer span.End()

    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id := vars["id"]
    configuration, err := h.r.FindById(ctx, id)

    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusNotFound, err.Error())
        return
    }

    err = json.NewEncoder(w).Encode(configuration)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
    } else {
        span.SetStatus(codes.Ok, "")
    }
}

// Create creates a new configuration
// swagger:route POST /configurations configurations createConfiguration
//
// # Create a new configuration
//
// This endpoint creates a new configuration with the provided data.
//
// Responses:
//
//	200: body:Configuration
//	400: body:ErrorResponse
//	409: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationHandler) Create(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "ConfigurationHandler.Create")
    defer span.End()

    w.Header().Set("Content-Type", "application/json")

    var configuration dtos.CreateConfigurationDto
    _ = json.NewDecoder(r.Body).Decode(&configuration)

    if len(configuration.Parameters) < 1 {
        span.SetStatus(codes.Error, "you must specify at least one parameter")
        h.sendErrorResponse(w, http.StatusBadRequest, "you must specify at least one parameter")
        return
    }

    createdConfig, err := h.r.Create(ctx, mappers.ToConfiguration(&configuration))
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusConflict, err.Error())
        return
    }

    err = json.NewEncoder(w).Encode(createdConfig)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
    } else {
        span.SetStatus(codes.Ok, "")
    }
}

// Delete removes a configuration by ID
// swagger:route DELETE /configurations/{id} configurations deleteConfigurationById
//
// # Delete configuration by ID
//
// This endpoint deletes a specific configuration by its ID.
//
// Parameters:
//   - name: id
//     in: path
//     type: string
//     required: true
//     description: The ID of the configuration
//
// Responses:
//
//	204: body:NoContentResponse
//	404: body:ErrorResponse
func (h ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "ConfigurationHandler.Delete")
    defer span.End()

    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id := vars["id"]
    err := h.r.DeleteById(ctx, id)

    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusNotFound, err.Error())
        return
    }

    span.SetStatus(codes.Ok, "")
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(models.NoContentResponse{})
}

// DeleteByNameAndVersion removes a configuration by name and version
// swagger:route DELETE /configuration/{name}/{version} configurations deleteConfigurationByNameAndVersion
//
// # Delete configuration by name and version
//
// This endpoint deletes a specific configuration by its name and version.
//
// Parameters:
//   - name: name
//     in: path
//     type: string
//     required: true
//     description: The name of the configuration
//   - name: version
//     in: path
//     type: string
//     required: true
//     description: The version of the configuration
//
// Responses:
//
//	204: body:NoContentResponse
//	404: body:ErrorResponse
func (h ConfigurationHandler) DeleteByNameAndVersion(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "ConfigurationHandler.DeleteByNameAndVersion")
    defer span.End()

    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    name := params["name"]
    version := params["version"]

    err := h.r.DeleteByNameAndVersion(ctx, name, version)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusNotFound, err.Error())
        return
    }

    span.SetStatus(codes.Ok, "")
    w.WriteHeader(http.StatusNoContent)
    json.NewEncoder(w).Encode(models.NoContentResponse{})
}

// FindByNameAndVersion retrieves a configuration by name and version
// swagger:route GET /configuration/{name}/{version} configurations getConfigurationByNameAndVersion
//
// # Get configuration by name and version
//
// This endpoint retrieves a specific configuration by its name and version.
//
// Parameters:
//   - name: name
//     in: path
//     type: string
//     required: true
//     description: The name of the configuration
//   - name: version
//     in: path
//     type: string
//     required: true
//     description: The version of the configuration
//
// Responses:
//
//	200: body:Configuration
//	404: body:ErrorResponse
//	500: body:ErrorResponse
func (h ConfigurationHandler) FindByNameAndVersion(w http.ResponseWriter, r *http.Request) {
    ctx, span := h.Tracer.Start(r.Context(), "ConfigurationHandler.FindByNameAndVersion")
    defer span.End()

    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    name := params["name"]
    version := params["version"]

    configuration, err := h.r.FindByNameAndVersion(ctx, name, version)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusNotFound, err.Error())
        return
    }

    err = json.NewEncoder(w).Encode(configuration)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        h.sendErrorResponse(w, http.StatusInternalServerError, err.Error())
    } else {
        span.SetStatus(codes.Ok, "")
    }
}