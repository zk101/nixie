package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

// Connect establishes a connection to consul based on the Config provided
func (c *Client) Connect() error {
	conf := consulapi.DefaultConfig()
	conf.Address = c.config.Address
	conf.Scheme = c.config.Scheme
	conf.Datacenter = c.config.Datacenter
	conf.Token = c.config.Token

	var err error
	c.conn, err = consulapi.NewClient(conf)
	if err != nil {
		return err
	}

	return nil
}

// EOF
