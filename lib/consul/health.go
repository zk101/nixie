package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

// HealthService checks a services health
func (c *Client) HealthService(service, tag string, passingOnly bool, q *consulapi.QueryOptions) ([]*consulapi.ServiceEntry, *consulapi.QueryMeta, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return nil, nil, err
		}
	}

	serviceEntries, queryMeta, err := c.conn.Health().Service(service, tag, passingOnly, q)
	if err != nil {
		return nil, nil, err
	}

	return serviceEntries, queryMeta, nil
}

// EOF
