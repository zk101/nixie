package couchbase

import (
	"github.com/couchbase/gocb"
)

// Client is an operation struct
type Client struct {
	config *Config
	conn   *gocb.Bucket
}

// NewClient returns a pointer to a new Client instance
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
