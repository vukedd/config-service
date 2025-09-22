package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/vukedd/config-service/dtos"
	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
	"net/http"
)

type ConfigurationGroupHandler struct {
	Repository              *repositories.ConfigurationGroupRepository
	ConfigurationRepository *repositories.ConfigurationRepository
}

func NewConfigurationGroupHandler(repository *repositories.ConfigurationGroupRepository, configurationRepository *repositories.ConfigurationRepository) *ConfigurationGroupHandler {
	return &ConfigurationGroupHandler{Repository: repository, ConfigurationRepository: configurationRepository}
}

func (Handler ConfigurationGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	configurationGroups := Handler.Repository.FindAll()
	err := json.NewEncoder(w).Encode(configurationGroups)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	return
}

func (Handler ConfigurationGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]

	configurationGroup, err := Handler.Repository.FindById(configurationGroupId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = json.NewEncoder(w).Encode(configurationGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	return
}

func (Handler ConfigurationGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]
	err := Handler.Repository.Delete(configurationGroupId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}

func (Handler ConfigurationGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var configurationGroupRequest dtos.ConfigurationGroupDto
	_ = json.NewDecoder(r.Body).Decode(&configurationGroupRequest)

	if len(configurationGroupRequest.ConfigurationList) < 1 {
		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]string{"error": "you must define at least one configuration"}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	for _, configuration := range Handler.Repository.ConfigurationGroups {
		if configuration.Name == configurationGroupRequest.Name && configuration.Version == configurationGroupRequest.Version {
			w.WriteHeader(http.StatusConflict)
			errorResponse := map[string]string{"error": "configuration group already exists"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// I thought it was a good idea to leave the transformation from dto the model in the handler since we are going
	// to check if selected configurations for this configuration group exist by fetching data from the repository,
	// and by leaving it as it is, I am avoiding giving mapper classes data access :D SAME GOES FOR THE UPDATE METHOD
	configurationGroupConfigurationList := []*models.LabeledConfiguration{}

	for _, configurationItem := range configurationGroupRequest.ConfigurationList {
		found := false
		for _, configuration := range Handler.ConfigurationRepository.Configurations {
			if configurationItem.Id == configuration.Id {
				found = true
				configurationGroupConfigurationList = append(configurationGroupConfigurationList, &models.LabeledConfiguration{Id: "", Configuration: configuration, Labels: configurationItem.Labels})
			}
		}

		if found == false {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := map[string]string{"error": "configuration with the id " + configurationItem.Id + " does not exist"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	newConfigurationGroup := models.ConfigurationGroup{Id: "", Name: configurationGroupRequest.Name, Version: configurationGroupRequest.Version, Configurations: configurationGroupConfigurationList}
	configGroup, err := Handler.Repository.Create(&newConfigurationGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	err = json.NewEncoder(w).Encode(configGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	return
}

func (Handler ConfigurationGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]

	var configGroupData dtos.ConfigurationGroupDto
	json.NewDecoder(r.Body).Decode(&configGroupData)

	if len(configGroupData.ConfigurationList) < 1 {
		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]string{"error": "you must define at least one configuration"}
		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	configurationGroupConfigurationList := []*models.LabeledConfiguration{}

	for _, configurationItem := range configGroupData.ConfigurationList {
		found := false
		for _, configuration := range Handler.ConfigurationRepository.Configurations {
			if configurationItem.Id == configuration.Id {
				found = true
				configurationGroupConfigurationList = append(configurationGroupConfigurationList, &models.LabeledConfiguration{Id: "", Configuration: configuration, Labels: configurationItem.Labels})
			}
		}

		if found == false {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := map[string]string{"error": "configuration with the id " + configurationItem.Id + " does not exist"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	updateConfigurationGroup := models.ConfigurationGroup{Id: configurationGroupId, Name: configGroupData.Name, Version: configGroupData.Version, Configurations: configurationGroupConfigurationList}
	err := Handler.Repository.Update(configurationGroupId, &updateConfigurationGroup)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	return
}

func (Handler ConfigurationGroupHandler) FindByIdToDto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	configurationGroupId := params["id"]

	configurationGroup, err := Handler.Repository.FindById(configurationGroupId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	configurationsTransformedToDto := []*dtos.ConfigurationGroupConfigurationDto{}
	for _, labeledConfiguration := range configurationGroup.Configurations {
		configurationsTransformedToDto = append(configurationsTransformedToDto, &dtos.ConfigurationGroupConfigurationDto{Id: labeledConfiguration.Configuration.Id, Labels: labeledConfiguration.Labels})
	}

	configurationGroupDto := dtos.ConfigurationGroupDto{Name: configurationGroup.Name, Version: configurationGroup.Version, ConfigurationList: configurationsTransformedToDto}
	err = json.NewEncoder(w).Encode(configurationGroupDto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := map[string]string{"error": err.Error()}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	return
}
