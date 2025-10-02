package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/dtos"
	"github.com/vukedd/config-service/mappers"
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

func (Handler ConfigurationHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	configurations := Handler.repository.FindAll()

	err := json.NewEncoder(w).Encode(configurations)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}
	return
}

func (Handler ConfigurationHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	configuration, err := Handler.repository.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	err = json.NewEncoder(w).Encode(configuration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	return
}

func (Handler ConfigurationHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var configuration dtos.CreateConfigurationDto
	_ = json.NewDecoder(r.Body).Decode(&configuration)

	if len(configuration.Parameters) < 1 {
		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]string{"error": "you must specify at least one parameter"}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	createdConfig, err := Handler.repository.Create(*(mappers.ToConfiguration(&configuration)))
	if err != nil {
		w.WriteHeader(http.StatusConflict)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	err = json.NewEncoder(w).Encode(createdConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	return
}

func (Handler ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	err := Handler.repository.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (Handler ConfigurationHandler) DeleteByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	name := params["name"]
	version := params["version"]

	err := Handler.repository.DeleteByNameAndVersion(name, version)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (Handler ConfigurationHandler) FindByNameAndVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	name := params["name"]
	version := params["version"]

	configuration, err := Handler.repository.FindByNameAndVersion(name, version)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = json.NewEncoder(w).Encode(configuration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	return
}
