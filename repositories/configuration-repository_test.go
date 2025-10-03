package repositories_test

import (
	"testing"

	. "github.com/franela/goblin"

	"github.com/vukedd/config-service/models"
	"github.com/vukedd/config-service/repositories"
)

func assertConfigurationsEqual(g *G, config1 *models.Configuration, config2 *models.Configuration) {
	if config1.Name != config2.Name {
		g.Errorf("Expected name to be %s, got %s", config2.Name, config1.Name)
	}

	if config1.Version != config2.Version {
		g.Errorf("Expected version to be %s, got %s", config2.Version, config1.Version)
	}

	if config1.Parameters["db_url"] != config2.Parameters["db_url"] {
		g.Errorf("Expected db_url to be %s, got %s", config2.Parameters["db_url"], config1.Parameters["db_url"])
	}
}

func TestFindByNameAndVersion(t *testing.T) {
	g := Goblin(t)
	c := createConsul(g)
	repo := repositories.NewConfigurationRepository(c)

	g.Describe("FindByNameAndVersion", func() {
		g.AfterEach(func() {
			_ = repo.DeleteByNameAndVersion("golang-test-1", "1.0.0")
			_ = repo.DeleteByNameAndVersion("golang-test-1", "1.1.0")
		})

		g.It("should return a configuration", func() {
			_, err := repo.Create(&models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}

			config, err := repo.FindByNameAndVersion("golang-test-1", "1.0.0")
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}

			assertConfigurationsEqual(g, config, &models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
		})

		g.It("should search by name and version", func() {
			_, err := repo.Create(&models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}
			_, err = repo.Create(&models.Configuration{Name: "golang-test-1", Version: "1.1.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}

			config1, err := repo.FindByNameAndVersion("golang-test-1", "1.1.0")
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}

			config2, err := repo.FindByNameAndVersion("golang-test-1", "1.0.0")
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}

			assertConfigurationsEqual(g, config1, &models.Configuration{Name: "golang-test-1", Version: "1.1.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			assertConfigurationsEqual(g, config2, &models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
		})

		g.It("should return error if configuration not found", func() {
			_, err := repo.FindByNameAndVersion("golang-test-1", "1.0.0")
			if err == nil {
				g.Errorf("Expected error, got nil")
			}

			if err.Error() != "configuration not found" {
				g.Errorf("Expected error to be 'configuration not found', got %s", err.Error())
			}
		})
	})
}

func TestCreate(t *testing.T) {
	g := Goblin(t)
	c := createConsul(g)
	repo := repositories.NewConfigurationRepository(c)

	g.Describe("Create", func() {
		g.AfterEach(func() {
			_ = repo.DeleteByNameAndVersion("golang-test-1", "1.0.0")
		})

		g.It("should create a configuration", func() {
			config, err := repo.Create(&models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			if err != nil {
				g.Errorf("Expected no error, got %s", err.Error())
			}

			assertConfigurationsEqual(g, config, &models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
		})

		g.It("should return error if configuration already exists", func() {
			repo.Create(&models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			_, err := repo.Create(&models.Configuration{Name: "golang-test-1", Version: "1.0.0", Parameters: map[string]string{"db_url": "localhost:1234/db"}})
			if err == nil {
				g.Errorf("Expected error, got nil")
			}

			if err.Error() != "configuration already exists" {
				g.Errorf("Expected error to be 'configuration already exists', got %s", err.Error())
			}
		})
	})
}
