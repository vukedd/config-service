package repositories_test

import (
	"context"
	"testing"

	. "github.com/franela/goblin"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
)

func assertConfigGroupEqual(g *G, configGroup1 *models.ConfigurationGroup, configGroup2 *models.ConfigurationGroup) {
    if configGroup1 == nil {
        g.Errorf("Expected configGroup1 to not be nil")
        return
    }

    if configGroup2 == nil {
        g.Errorf("Expected configGroup2 to not be nil")
        return
    }

    if configGroup1.Name != configGroup2.Name {
        g.Errorf("Expected Name to be %s, got %s", configGroup2.Name, configGroup1.Name)
    }

    if configGroup1.Version != configGroup2.Version {
        g.Errorf("Expected Version to be %s, got %s", configGroup2.Version, configGroup1.Version)
    }

    if len(configGroup1.Configurations) != len(configGroup2.Configurations) {
        g.Errorf("Expected %d configurations, got %d", len(configGroup2.Configurations), len(configGroup1.Configurations))
        return
    }

    for i, config1 := range configGroup1.Configurations {
        if i < len(configGroup2.Configurations) {
            config2 := configGroup2.Configurations[i]
            if config1.Configuration != nil && config2.Configuration != nil {
                assertConfigurationsEqual(g, config1.Configuration, config2.Configuration)
            }
            for key, value := range config1.Labels {
                if config2.Labels[key] != value {
                    g.Errorf("Expected label %s to be %s, got %s", key, value, config2.Labels[key])
                }
            }
        }
    }
}

func TestConfigurationGroupRepository(t *testing.T) {
    g := Goblin(t)
    c := createConsul(g)
    
    // Create tracers for testing
    tracer := noop.NewTracerProvider().Tracer("test")
    gr := repositories.NewConfigurationGroupRepository(c, tracer)
    cr := repositories.NewConfigurationRepository(c, tracer)

    g.Describe("ConfigurationGroupRepository", func() {
        g.Describe("FindById", func() {
            g.AfterEach(func() {
                ctx := context.Background()
                _ = gr.DeleteByNameAndVersion(ctx, "golang-test-1", "1.0.0")
                _ = cr.DeleteByNameAndVersion(ctx, "golang-test-config", "1.0.0")
            })

            g.It("should return a configuration group", func() {
                ctx := context.Background()
                
                // First create a configuration that the group will reference
                config, err := cr.Create(ctx, &models.Configuration{
                    Name:       "golang-test-config",
                    Version:    "1.0.0",
                    Parameters: map[string]string{"db_url": "localhost:1234/db"},
                })
                if err != nil {
                    g.Errorf("Failed to create configuration: %s", err.Error())
                    return
                }

                // Now create the configuration group
                group, err := gr.Create(ctx, &models.ConfigurationGroup{
                    Name:    "golang-test-1",
                    Version: "1.0.0",
                    Configurations: []*models.LabeledConfiguration{
                        {
                            Configuration: config,
                            Labels:        map[string]string{"env": "production", "region": "eu-central"},
                        },
                    },
                })

                if err != nil {
                    g.Errorf("Expected no error, got %s", err.Error())
                    return
                }

                configGroup, err := gr.FindById(ctx, group.Id)
                if err != nil {
                    g.Errorf("Expected no error, got %s", err.Error())
                    return
                }

                expectedGroup := &models.ConfigurationGroup{
                    Name:    "golang-test-1",
                    Version: "1.0.0",
                    Configurations: []*models.LabeledConfiguration{
                        {
                            Configuration: config,
                            Labels:        map[string]string{"env": "production", "region": "eu-central"},
                        },
                    },
                }

                assertConfigGroupEqual(g, configGroup, expectedGroup)
            })
        })
    })
}