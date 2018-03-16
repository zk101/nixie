package consul

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

// GetKV returns a KeyValue from Consul
func (c *Client) GetKV(key string) (string, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}

	kv := c.conn.KV()
	pair, _, err := kv.Get(key, nil)

	if err != nil {
		return "", err
	}

	if pair == nil {
		return "", nil
	}

	return string(pair.Value), nil
}

// Get is a wrapper function for GetKV to statisfy config kvstore
func (c *Client) Get(key string) (string, error) {
	return c.GetKV(strings.ToLower(key))
}

// PutKV returns a KeyValue from Consul
func (c *Client) PutKV(key string, value string) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	kv := c.conn.KV()
	p := &consulapi.KVPair{Key: key, Value: []byte(value)}

	if _, err := kv.Put(p, nil); err != nil {
		return err
	}

	return nil
}

// EOF
