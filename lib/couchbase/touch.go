package couchbase

import (
	"github.com/couchbase/gocb"
)

// Touch updates expiry on a document
func (c *Client) Touch(key string, cas gocb.Cas, expiry uint32) (gocb.Cas, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return 0, err
		}
	}

	cas, err := c.conn.Touch(key, cas, expiry)
	if err != nil {
		return 0, err
	}

	return cas, nil
}

// EOF
