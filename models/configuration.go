package models

type Configuration struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Version    string            `json:"version"`
	Parameters map[string]string `json:"parameters"`
}

type ConfigurationGroup struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	Version        string           `json:"version"`
	Configurations []*Configuration `json:"configurations"`
}
