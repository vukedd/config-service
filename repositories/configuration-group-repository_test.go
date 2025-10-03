package repositories_test

import (
	"testing"

	. "github.com/franela/goblin"
	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
)

func assertConfigGroupEqual(t *G, configGroup1 *models.ConfigurationGroup, configGroup2 *models.ConfigurationGroup) {
	if configGroup1.Name != configGroup2.Name {
		t.Errorf("Expected Name to be %s, got %s", configGroup2.Name, configGroup1.Name)
	}
	
	if configGroup1.Version != configGroup2.Version {
		t.Errorf("Expected Version to be %s, got %s", configGroup2.Version, configGroup1.Version)
	}

	for _, config1 := range configGroup1.Configurations {
		for _, config2 := range configGroup2.Configurations {
			assertConfigurationsEqual(t, config1.Configuration, config2.Configuration)
			for key, value := range config1.Labels {
				if config2.Labels[key] != value {
					t.Errorf("Expected label %s to be %s, got %s", key, value, config2.Labels[key])
				}
			}
		}
	}
}

func TestConfigurationGroupRepository(t *testing.T) {
	g := Goblin(t)
	g.Describe("ConfigurationGroupRepository", func() {
		g.Describe("FindById", func() {
			g.It("should return a configuration group", func() {
				groupRepo := repositories.NewConfigurationGroupRepository()

				group, err := groupRepo.Create(&models.ConfigurationGroup{
					Name:    "golang-test-1",
					Version: "1.0.0",
					Configurations: []*models.LabeledConfiguration{
						{
							Configuration: &models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}},
							Labels:        map[string]string{"env": "production", "region": "eu-central"},
						},
					},
				})

				if err != nil {
					g.Errorf("Expected no error, got %s", err.Error())
				}

				configGroup, err := groupRepo.FindById(group.Id)
				if err != nil {
					g.Errorf("Expected no error, got %s", err.Error())
				}

				assertConfigGroupEqual(g, configGroup, &models.ConfigurationGroup{
					Name:    "golang-test-1",
					Version: "1.0.0",
					Configurations: []*models.LabeledConfiguration{
						{
							Configuration: &models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}},
							Labels:        map[string]string{"env": "production", "region": "eu-central"},
						},
					},
				})
			})
		})
	})
}
