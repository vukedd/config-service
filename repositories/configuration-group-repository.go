package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/vukedd/config-service/models"
)

const consulGroupsKey = "groups"

var (
	ErrConfigurationGroupNotFound = errors.New("configuration group not found")
	ErrConfigurationGroupExists   = errors.New("configuration group already exists")
	ErrDuplicateLabel             = errors.New("duplicate label")
)

type ConfigurationGroupRepository struct {
	*ConsulClient
}

func (c *ConfigurationGroupRepository) kvKeyFromConfigurationGroup(configGroup *models.ConfigurationGroup) string {
	return c.kvKey(consulGroupsKey, configGroup.Name, configGroup.Version)
}

func NewConfigurationGroupRepository(consulClient *api.Client) *ConfigurationGroupRepository {
	return &ConfigurationGroupRepository{
		ConsulClient: &ConsulClient{consulClient},
	}
}

func (r *ConfigurationGroupRepository) FindAll() ([]*models.ConfigurationGroup, error) {
	kv := r.consul.KV()
	pairs, _, err := kv.List(fmt.Sprintf("%s/", consulGroupsKey), nil)
	if err != nil {
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

	return configs, nil
}

func (r *ConfigurationGroupRepository) FindById(id string) (*models.ConfigurationGroup, error) {
	groups, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Id == id {
			return group, nil
		}
	}
	return nil, ErrConfigurationGroupNotFound
}

// extractLabels converts list argument of (key:value;key2:value) to map[string]string
func extractLabels(list string) (map[string]string, error) {
	split := strings.Split(list, ";")
	values := make(map[string]string)

	for _, s := range split {
		kv := strings.Split(s, ":")
		if _, ok := values[kv[0]]; ok {
			return nil, fmt.Errorf("%w: %s", ErrDuplicateLabel, kv[0])
		}

		values[kv[0]] = kv[1]
	}

	return values, nil
}

func (r *ConfigurationGroupRepository) FindByLabel(list string) ([]*models.ConfigurationGroup, error) {
	groups, err := r.FindAll()
	if err != nil {
		return nil, err
	}

	labels, err := extractLabels(list)
	if err != nil {
		return nil, err
	}

	cg := []*models.ConfigurationGroup{}
	for _, group := range groups {
		for _, lc := range group.Configurations {
			found := true
			for k, v := range labels {
				if lc.Labels[k] != v {
					found = false
					break
				}
			}

			if found {
				cg = append(cg, group)
			}
		}
	}

	return cg, nil
}

func (r *ConfigurationGroupRepository) Create(g *models.ConfigurationGroup) (*models.ConfigurationGroup, error) {
	kv := r.consul.KV()
	key := r.kvKeyFromConfigurationGroup(g)

	existing, _, _ := kv.Get(key, nil)
	if existing != nil {
		return nil, ErrConfigurationGroupExists
	}

	g.Id = uuid.New().String()
	data, _ := json.Marshal(g)

	p := &api.KVPair{Key: key, Value: data}
	_, err := kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (r *ConfigurationGroupRepository) DeleteByNameAndVersion(name string, version string) error {
	kv := r.consul.KV()
	key := r.kvKey(consulGroupsKey, name, version)
	_, err := kv.Delete(key, nil)
	return err
}

func (r *ConfigurationGroupRepository) DeleteById(id string) error {
	groups, err := r.FindAll()
	if err != nil {
		return err
	}

	kv := r.consul.KV()
	for _, g := range groups {
		if g.Id == id {
			key := r.kvKeyFromConfigurationGroup(g)
			_, err := kv.Delete(key, nil)
			return err
		}
	}

	return ErrConfigurationNotFound
}

// Update updates an existing configuration group
// with new configurations from `cg`
func (r *ConfigurationGroupRepository) Update(Id string, cg *models.ConfigurationGroup) error {
	g, err := r.FindById(Id)
	if err != nil {
		return err
	}

	g.Configurations = cg.Configurations
	return nil
}

func (r *ConfigurationGroupRepository) FindByNameAndVersion(name string, version string) (*models.ConfigurationGroup, error) {
	kv := r.consul.KV()
	key := r.kvKey(consulGroupsKey, name, version)

	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	if pair == nil {
		return nil, ErrConfigurationGroupNotFound
	}

	var configGroup models.ConfigurationGroup
	if err := json.Unmarshal(pair.Value, &configGroup); err != nil {
		return nil, err
	}

	return &configGroup, nil
}
