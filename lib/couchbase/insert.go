package couchbase

import (
	"github.com/couchbase/gocb"
)

// Insert places a document into couchbase
func (c *Client) Insert(key string, value interface{}, expiry uint32) (gocb.Cas, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return 0, err
		}
	}

	cas, err := c.conn.Insert(key, value, expiry)
	if err != nil {
		return 0, err
	}

	return cas, nil
}

// EOF
