package dtos

type ConfigurationGroupDto struct {
	Name              string                                `json:"name"`
	Version           string                                `json:"version"`
	ConfigurationList []*ConfigurationGroupConfigurationDto `json:"configuration_list"`
}
