package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

// RegisterService attempts to register a service in consul
func (c *Client) RegisterService(service *consulapi.AgentServiceRegistration) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	return c.conn.Agent().ServiceRegister(service)
}

// UnregisterService attempts to un-register a service in consul
func (c *Client) UnregisterService(serviceID string) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	return c.conn.Agent().ServiceDeregister(serviceID)
}

// EOF
