package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/vukedd/config-service/models"
	"strconv"
)

type ConfigurationGroupRepository struct {
	ConfigurationGroups []*models.ConfigurationGroup
}

func NewConfigurationGroupRepository() *ConfigurationGroupRepository {
	configurationGroupRepository := ConfigurationGroupRepository{}
	for i := 0; i < 4; i++ {
		configurationGroupRepository.ConfigurationGroups = append(configurationGroupRepository.ConfigurationGroups, &models.ConfigurationGroup{
			uuid.New().String(),
			"Configuration group" + strconv.Itoa(i),
			"1.0." + strconv.Itoa(i),
			[]*models.LabeledConfiguration{
				&(models.LabeledConfiguration{
					Id: uuid.New().String(),
					Configuration: &models.Configuration{
						Id:         uuid.New().String(),
						Name:       "Config " + strconv.Itoa(i+1),
						Version:    "1.0." + strconv.Itoa(i+1),
						Parameters: map[string]string{"db_url": "db:330" + strconv.Itoa(i) + "/db"},
					},
					Labels: map[string]string{"env": "production", "region": "eu-central"},
				}),
				&(models.LabeledConfiguration{
					Id: uuid.New().String(),
					Configuration: &models.Configuration{
						Id:         uuid.New().String(),
						Name:       "Config " + strconv.Itoa(i+1),
						Version:    "1.0." + strconv.Itoa(i+1),
						Parameters: map[string]string{"db_url": "db:330" + strconv.Itoa(i) + "/db"},
					},
					Labels: map[string]string{"env": "production", "region": "eu-central"},
				}),
				&(models.LabeledConfiguration{
					Id: uuid.New().String(),
					Configuration: &models.Configuration{
						Id:         uuid.New().String(),
						Name:       "Config " + strconv.Itoa(i+1),
						Version:    "1.0." + strconv.Itoa(i+1),
						Parameters: map[string]string{"db_url": "db:330" + strconv.Itoa(i) + "/db"},
					},
					Labels: map[string]string{"env": "production", "region": "eu-central"},
				}),
			},
		})
	}

	return &configurationGroupRepository
}

func (Repository *ConfigurationGroupRepository) FindAll() []*models.ConfigurationGroup {
	return Repository.ConfigurationGroups
}

func (Repository *ConfigurationGroupRepository) FindById(Id string) (*models.ConfigurationGroup, error) {
	for _, configurationGroup := range Repository.ConfigurationGroups {
		if configurationGroup.Id == Id {
			return configurationGroup, nil
		}
	}
	return &models.ConfigurationGroup{"", "", "", []*models.LabeledConfiguration{}}, errors.New("configuration group not found")
}

func (Repository *ConfigurationGroupRepository) Delete(Id string) error {
	for _, configurationGroup := range Repository.ConfigurationGroups {
		if configurationGroup.Id == Id {
			Repository.ConfigurationGroups = append(Repository.ConfigurationGroups[:0], Repository.ConfigurationGroups[1:]...)
			return nil
		}
	}
	return errors.New("configuration group not found")
}

func (Repository *ConfigurationGroupRepository) Create(ConfigurationGroup *models.ConfigurationGroup) (*models.ConfigurationGroup, error) {
	ConfigurationGroup.Id = uuid.New().String()
	for _, configuration := range ConfigurationGroup.Configurations {
		configuration.Id = uuid.New().String()
	}

	Repository.ConfigurationGroups = append(Repository.ConfigurationGroups, ConfigurationGroup)
	return ConfigurationGroup, nil
}

func (Repository *ConfigurationGroupRepository) Update(Id string, ConfigurationGroup *models.ConfigurationGroup) error {
	targetIndex := -1
	for i, configuration := range Repository.ConfigurationGroups {
		if configuration.Id == Id {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return errors.New("configuration group not found")
	}

	Repository.ConfigurationGroups[targetIndex].Configurations = ConfigurationGroup.Configurations
	return nil
}
