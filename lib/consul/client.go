package consul

import (
	consulapi "github.com/hashicorp/consul/api"
)

// Client holds consul details
type Client struct {
	config *Config
	conn   *consulapi.Client
}

// NewClient setups a new Client struct and returns a pointer to it
func NewClient(conf *Config) *Client {
	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	return &Client{
		config: conf,
	}
}

// EOF
