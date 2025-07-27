package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
	"net/http"
)

type ConfigurationHandler struct {
	repository *repositories.Repository
}

func NewConfigurationHandler(repository *repositories.Repository) *ConfigurationHandler {
	return &ConfigurationHandler{
		repository: repository,
	}
}

func (handler ConfigurationHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	configurations := handler.repository.FindAll()

	err := json.NewEncoder(w).Encode(configurations)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}
	return
}

func (handler ConfigurationHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	configuration, err := handler.repository.FindById(id)

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

func (handler ConfigurationHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var configuration models.Configuration
	_ = json.NewDecoder(r.Body).Decode(&configuration)

	if len(configuration.Parameters) < 1 {
		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]string{"error": "you must specify at least one parameter"}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	createdConfig, err := handler.repository.Create(configuration)
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

func (handler ConfigurationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	err := handler.repository.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
