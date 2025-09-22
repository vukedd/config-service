package mappers

import (
	"github.com/vukedd/config-service/dtos"
	"github.com/vukedd/config-service/models"
)

func ToConfiguration(configurationDto *dtos.CreateConfigurationDto) *models.Configuration {
	return &models.Configuration{
		Id:         "",
		Name:       configurationDto.Name,
		Version:    configurationDto.Version,
		Parameters: configurationDto.Parameters,
	}
}
