package repositories

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	consul *api.Client
}

// kvKey is a helper function to generate a KV key.
func (c *ConsulClient) kvKey(group, name, version string) string {
	return fmt.Sprintf("%s/%s/%s", group, name, version)
}
