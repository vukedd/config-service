package dtos

type ConfigurationGroupConfigurationDto struct {
	Id     string            `json:"id"`
	Labels map[string]string `json:"labels"`
}

type CreateConfigurationDto struct {
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Parameters map[string]string `json:"parameters"`
}
