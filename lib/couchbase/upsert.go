package couchbase

import (
	"github.com/couchbase/gocb"
)

// Upsert places a document into couchbase
func (c *Client) Upsert(key string, value interface{}, expiry uint32) (gocb.Cas, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return 0, err
		}
	}

	cas, err := c.conn.Upsert(key, value, expiry)
	if err != nil {
		return 0, err
	}

	return cas, nil
}

// EOF
