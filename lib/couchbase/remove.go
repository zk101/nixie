package couchbase

import (
	"github.com/couchbase/gocb"
)

// Remove deletes a document in couchbase
func (c *Client) Remove(key string, cas gocb.Cas) (gocb.Cas, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return 0, err
		}
	}

	cas, err := c.conn.Remove(key, cas)
	if err != nil {
		return 0, err
	}

	return cas, nil
}

// EOF
