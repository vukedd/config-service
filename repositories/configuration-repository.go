package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/models"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const consulConfigsKey = "configs"

var (
	ErrConfigurationNotFound = errors.New("configuration not found")
	ErrConfigurationExists   = errors.New("configuration already exists")
)

type ConfigurationRepository struct {
	*ConsulClient
	Tracer trace.Tracer
}

func NewConfigurationRepository(consulClient *api.Client, tracer trace.Tracer) *ConfigurationRepository {
	r := &ConfigurationRepository{
		ConsulClient: &ConsulClient{consulClient},
		Tracer:       tracer,
	}
	r.prepopulateData()
	return r
}

func (r *ConfigurationRepository) kvKeyFromConfiguration(config *models.Configuration) string {
	return r.kvKey(consulConfigsKey, config.Name, config.Version)
}

// Find all configurations
func (r *ConfigurationRepository) FindAll(ctx context.Context) ([]*models.Configuration, error) {
	_, span := r.Tracer.Start(ctx, "ConfigurationRepository.FindAll")
	defer span.End()

	kv := r.consul.KV()
	pairs, _, err := kv.List(fmt.Sprintf("%s/", consulConfigsKey), nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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

	span.SetStatus(codes.Ok, "")
	return configs, nil
}

func (r *ConfigurationRepository) FindById(ctx context.Context, id string) (*models.Configuration, error) {
	_, span := r.Tracer.Start(ctx, "ConfigurationRepository.FindById")
	defer span.End()

	configurations, err := r.FindAll(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	for _, configuration := range configurations {
		if configuration.Id == id {
			span.SetStatus(codes.Ok, "")
			return configuration, nil
		}
	}

	span.SetStatus(codes.Error, "configuration not found")
	return nil, ErrConfigurationNotFound
}

// Create a new configuration
func (r *ConfigurationRepository) Create(ctx context.Context, c *models.Configuration) (*models.Configuration, error) {
	_, span := r.Tracer.Start(ctx, "ConfigurationRepository.Create")
	defer span.End()

	kv := r.consul.KV()
	key := r.kvKeyFromConfiguration(c)

	existing, _, _ := kv.Get(key, nil)
	if existing != nil {
		span.SetStatus(codes.Error, "configuration already exists")
		return nil, ErrConfigurationExists
	}

	c.Id = uuid.New().String()
	data, _ := json.Marshal(c)

	p := &api.KVPair{Key: key, Value: data}
	_, err := kv.Put(p, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "")
	return c, nil
}

// Delete by name and version
func (r *ConfigurationRepository) DeleteByNameAndVersion(ctx context.Context, name, version string) error {
	_, span := r.Tracer.Start(ctx, "ConfigurationRepository.DeleteByNameAndVersion")
	defer span.End()

	kv := r.consul.KV()
	key := r.kvKey(consulConfigsKey, name, version)
	_, err := kv.Delete(key, nil)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

func (r *ConfigurationRepository) DeleteById(ctx context.Context, id string) error {
	_, span := r.Tracer.Start(ctx, "ConfigurationRepository.DeleteById")
	defer span.End()

	configs, err := r.FindAll(ctx)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	kv := r.consul.KV()
	for _, c := range configs {
		if c.Id == id {
			key := r.kvKeyFromConfiguration(c)
			_, err := kv.Delete(key, nil)
			if err != nil {
				span.SetStatus(codes.Error, err.Error())
			} else {
				span.SetStatus(codes.Ok, "")
			}
			return err
		}
	}

	span.SetStatus(codes.Error, "configuration not found")
	return ErrConfigurationNotFound
}

// Find by name and version
func (r *ConfigurationRepository) FindByNameAndVersion(ctx context.Context, name, version string) (*models.Configuration, error) {
	_, span := r.Tracer.Start(ctx, "ConfigurationRepository.FindByNameAndVersion")
	defer span.End()

	kv := r.consul.KV()
	key := r.kvKey(consulConfigsKey, name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	if pair == nil {
		span.SetStatus(codes.Error, "configuration not found")
		return nil, ErrConfigurationNotFound
	}

	var config models.Configuration
	if err := json.Unmarshal(pair.Value, &config); err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetStatus(codes.Ok, "")
	return &config, nil
}

func (r *ConfigurationRepository) prepopulateData() {
	// If data already exists or an error occurs, just skip it
	pairs, _, err := r.consul.KV().List(consulConfigsKey, nil)
	if err != nil || len(pairs) > 0 {
		return
	}

	ctx := context.Background()
	for i := 0; i < 4; i++ {
		config := models.Configuration{
			Id:      uuid.New().String(),
			Name:    "config-" + strconv.Itoa(i),
			Version: "1.0." + strconv.Itoa(i),
			Parameters: map[string]string{
				"db_url": "db:3306/db",
			},
		}

		_, err := r.Create(ctx, &config)
		if err != nil {
			return
		}
	}
}
