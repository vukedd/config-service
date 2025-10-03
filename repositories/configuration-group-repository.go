package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/models"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const consulGroupsKey = "groups"

var (
    ErrConfigurationGroupNotFound = errors.New("configuration group not found")
    ErrConfigurationGroupExists   = errors.New("configuration group already exists")
)

type ConfigurationGroupRepository struct {
    *ConsulClient
    Tracer trace.Tracer
}

func (c *ConfigurationGroupRepository) kvKeyFromConfigurationGroup(configGroup *models.ConfigurationGroup) string {
    return c.kvKey(consulGroupsKey, configGroup.Name, configGroup.Version)
}

func NewConfigurationGroupRepository(consulClient *api.Client, tracer trace.Tracer) *ConfigurationGroupRepository {
    return &ConfigurationGroupRepository{
        ConsulClient: &ConsulClient{consulClient},
        Tracer:       tracer,
    }
}

func (r *ConfigurationGroupRepository) FindAll(ctx context.Context) ([]*models.ConfigurationGroup, error) {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.FindAll")
    defer span.End()

    kv := r.consul.KV()
    pairs, _, err := kv.List(fmt.Sprintf("%s/", consulGroupsKey), nil)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    configs := []*models.ConfigurationGroup{}
    for _, pair := range pairs {
        var c models.ConfigurationGroup
        if err := json.Unmarshal(pair.Value, &c); err != nil {
            continue // skip malformed entries
        }

        configs = append(configs, &c)
    }

    span.SetStatus(codes.Ok, "")
    return configs, nil
}

func (r *ConfigurationGroupRepository) FindById(ctx context.Context, id string) (*models.ConfigurationGroup, error) {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.FindById")
    defer span.End()

    groups, err := r.FindAll(ctx)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    for _, group := range groups {
        if group.Id == id {
            span.SetStatus(codes.Ok, "")
            return group, nil
        }
    }
    
    span.SetStatus(codes.Error, "configuration group not found")
    return nil, ErrConfigurationGroupNotFound
}

func (r *ConfigurationGroupRepository) Create(ctx context.Context, g *models.ConfigurationGroup) (*models.ConfigurationGroup, error) {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.Create")
    defer span.End()

    kv := r.consul.KV()
    key := r.kvKeyFromConfigurationGroup(g)

    existing, _, _ := kv.Get(key, nil)
    if existing != nil {
        span.SetStatus(codes.Error, "configuration group already exists")
        return nil, ErrConfigurationGroupExists
    }

    g.Id = uuid.New().String()
    data, _ := json.Marshal(g)

    p := &api.KVPair{Key: key, Value: data}
    _, err := kv.Put(p, nil)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    span.SetStatus(codes.Ok, "")
    return g, nil
}

func (r *ConfigurationGroupRepository) DeleteByNameAndVersion(ctx context.Context, name string, version string) error {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.DeleteByNameAndVersion")
    defer span.End()

    kv := r.consul.KV()
    key := r.kvKey(consulGroupsKey, name, version)
    _, err := kv.Delete(key, nil)
    
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
    } else {
        span.SetStatus(codes.Ok, "")
    }
    
    return err
}

func (r *ConfigurationGroupRepository) DeleteById(ctx context.Context, id string) error {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.DeleteById")
    defer span.End()

    groups, err := r.FindAll(ctx)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    kv := r.consul.KV()
    for _, g := range groups {
        if g.Id == id {
            key := r.kvKeyFromConfigurationGroup(g)
            _, err := kv.Delete(key, nil)
            if err != nil {
                span.SetStatus(codes.Error, err.Error())
            } else {
                span.SetStatus(codes.Ok, "")
            }
            return err
        }
    }

    span.SetStatus(codes.Error, "configuration group not found")
    return ErrConfigurationGroupNotFound
}

// Update updates an existing configuration group
// with new configurations from `cg`
func (r *ConfigurationGroupRepository) Update(ctx context.Context, Id string, cg *models.ConfigurationGroup) error {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.Update")
    defer span.End()

    g, err := r.FindById(ctx, Id)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        return err
    }

    g.Configurations = cg.Configurations
    span.SetStatus(codes.Ok, "")
    return nil
}

func (r *ConfigurationGroupRepository) FindByNameAndVersion(ctx context.Context, name string, version string) (*models.ConfigurationGroup, error) {
    _, span := r.Tracer.Start(ctx, "ConfigurationGroupRepository.FindByNameAndVersion")
    defer span.End()

    kv := r.consul.KV()
    key := r.kvKey(consulGroupsKey, name, version)

    pair, _, err := kv.Get(key, nil)
    if err != nil {
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    if pair == nil {
        span.SetStatus(codes.Error, "configuration group not found")
        return nil, ErrConfigurationGroupNotFound
    }

    var configGroup models.ConfigurationGroup
    if err := json.Unmarshal(pair.Value, &configGroup); err != nil {
        span.SetStatus(codes.Error, err.Error())
        return nil, err
    }

    span.SetStatus(codes.Ok, "")
    return &configGroup, nil
}