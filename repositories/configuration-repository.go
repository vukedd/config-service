package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/models"
)

const consulConfigsKey = "configs"

var (
	ErrConfigurationNotFound = errors.New("configuration not found")
	ErrConfigurationExists   = errors.New("configuration already exists")
)

type ConfigurationRepository struct {
	*ConsulClient
}

func NewConfigurationRepository(consulClient *api.Client) *ConfigurationRepository {
	r := &ConfigurationRepository{
		ConsulClient: &ConsulClient{consulClient},
	}
	r.prepopulateData()
	return r
}

func (r *ConfigurationRepository) kvKeyFromConfiguration(config *models.Configuration) string {
	return r.kvKey(consulConfigsKey, config.Name, config.Version)
}

// Find all configurations
func (r *ConfigurationRepository) FindAll() ([]*models.Configuration, error) {
	kv := r.consul.KV()
	pairs, _, err := kv.List(fmt.Sprintf("%s/", consulConfigsKey), nil)
	if err != nil {
		return nil, err
	}

	configs := []*models.Configuration{}
	for _, pair := range pairs {
		var c models.Configuration
		if err := json.Unmarshal(pair.Value, &c); err != nil {
			continue // skip malformed entries
		}

		configs = append(configs, &c)
	}

	return configs, nil
}

func (r *ConfigurationRepository) FindById(id string) (*models.Configuration, error) {
	configurations, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, configuration := range configurations {
		if configuration.Id == id {
			return configuration, nil
		}
	}
	return nil, ErrConfigurationNotFound
}

// Create a new configuration
func (r *ConfigurationRepository) Create(c *models.Configuration) (*models.Configuration, error) {
	kv := r.consul.KV()
	key := r.kvKeyFromConfiguration(c)

	existing, _, _ := kv.Get(key, nil)
	if existing != nil {
		return nil, ErrConfigurationExists
	}

	c.Id = uuid.New().String()
	data, _ := json.Marshal(c)

	p := &api.KVPair{Key: key, Value: data}
	_, err := kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Delete by name and version
func (r *ConfigurationRepository) DeleteByNameAndVersion(name, version string) error {
	kv := r.consul.KV()
	key := r.kvKey(consulConfigsKey, name, version)
	_, err := kv.Delete(key, nil)
	return err
}

func (r *ConfigurationRepository) DeleteById(id string) error {
	configs, err := r.FindAll()
	if err != nil {
		return err
	}

	kv := r.consul.KV()
	for _, c := range configs {
		if c.Id == id {
			key := r.kvKeyFromConfiguration(c)
			_, err := kv.Delete(key, nil)
			return err
		}
	}

	return ErrConfigurationNotFound
}

// Find by name and version
func (r *ConfigurationRepository) FindByNameAndVersion(name, version string) (*models.Configuration, error) {
	kv := r.consul.KV()
	key := r.kvKey(consulConfigsKey, name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, ErrConfigurationNotFound
	}

	var config models.Configuration
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigurationRepository) prepopulateData() {
	// If data already exists or an error occurs, just skip it
	pairs, _, err := r.consul.KV().List(consulConfigsKey, nil)
	if err != nil || len(pairs) > 0 {
		return
	}

	for i := 0; i < 4; i++ {
		config := models.Configuration{
			Id:      uuid.New().String(),
			Name:    "config-" + strconv.Itoa(i),
			Version: "1.0." + strconv.Itoa(i),
			Parameters: map[string]string{
				"db_url": "db:3306/db",
			},
		}

		_, err := r.Create(&config)
		if err != nil {
			return
		}
	}
}
