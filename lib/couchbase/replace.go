package couchbase

import (
	"github.com/couchbase/gocb"
)

// Replace deletes a document and adds it again in couchbase
func (c *Client) Replace(key string, value interface{}, cas gocb.Cas, expiry uint32) (gocb.Cas, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return 0, err
		}
	}

	cas, err := c.conn.Replace(key, value, cas, expiry)
	if err != nil {
		return 0, err
	}

	return cas, nil
}

// EOF
