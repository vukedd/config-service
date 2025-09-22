package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/vukedd/config-service/models"
	"strconv"
)

type ConfigurationRepository struct {
	Configurations []*models.Configuration
}

func NewRepository() *ConfigurationRepository {
	repository := ConfigurationRepository{}
	for i := 0; i < 4; i++ {
		repository.Configurations = append(repository.Configurations, &(models.Configuration{Id: uuid.New().String(), Name: "Config " + strconv.Itoa(i), Version: "1.0." + strconv.Itoa(i), Parameters: map[string]string{"db_url": "db:3306/db"}}))
	}

	return &repository
}

func (repository *ConfigurationRepository) FindAll() []*models.Configuration {
	return repository.Configurations
}

func (repository *ConfigurationRepository) FindById(id string) (*models.Configuration, error) {
	for _, configuration := range repository.Configurations {
		if configuration.Id == id {
			return configuration, nil
		}
	}
	return &models.Configuration{Id: "", Name: "", Version: "", Parameters: map[string]string{}}, errors.New("configuration not found")
}

func (Repository *ConfigurationRepository) Create(Configuration models.Configuration) (*models.Configuration, error) {
	for _, configuration := range Repository.Configurations {
		if configuration.Name == Configuration.Name && configuration.Version == Configuration.Version {
			return &Configuration, errors.New("configuration already exists")
		}
	}

	Configuration.Id = uuid.New().String()
	Repository.Configurations = append(Repository.Configurations, &Configuration)
	return &Configuration, nil
}

func (Repository *ConfigurationRepository) Delete(id string) error {
	for i, configuration := range Repository.Configurations {
		if configuration.Id == id {
			Repository.Configurations = append(Repository.Configurations[:i], Repository.Configurations[i+1:]...)
			return nil
		}
	}

	return errors.New("configuration not found")
}
