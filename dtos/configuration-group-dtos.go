package dtos

type CreateConfigurationGroupRequest struct {
	Name              string                                      `json:"name"`
	Version           string                                      `json:"version"`
	ConfigurationList []*CreateConfigurationGroupConfigurationDto `json:"configuration_list"`
}
