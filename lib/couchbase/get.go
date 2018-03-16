package couchbase

import (
	"github.com/couchbase/gocb"
)

// Get retrieves a document from a Bucket configuration
func (c *Client) Get(key string, valuePtr interface{}) (gocb.Cas, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return 0, err
		}
	}

	cas, err := c.conn.Get(key, valuePtr)
	if err != nil {
		return 0, err
	}

	return cas, nil
}

// EOF
