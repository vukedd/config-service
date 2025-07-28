package dtos

type ConfigurationGroupConfigurationDto struct {
	Id     string            `json:"id"`
	Labels map[string]string `json:"labels"`
}
