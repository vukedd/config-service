package repositories_test

import (
	. "github.com/franela/goblin"
	"github.com/hashicorp/consul/api"
)

func createConsul(g *G) *api.Client {
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := api.NewClient(config)
	if err != nil {
		g.Errorf("Expected no error, got %s", err.Error())
	}
	return client
}
