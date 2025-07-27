package dtos

type CreateConfigurationGroupConfigurationDto struct {
	Id     string            `json:"id"`
	Labels map[string]string `json:"labels"`
}
