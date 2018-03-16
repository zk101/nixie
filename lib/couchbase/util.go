package couchbase

import (
	"github.com/couchbase/gocb"
)

// Connect opens a connection to a couchbase bucket using Config
func (c *Client) Connect() error {
	if c.conn != nil {
		c.Close()
	}

	cluster, err := gocb.Connect(c.config.Cluster)
	if err != nil {
		return err
	}

	auth := gocb.PasswordAuthenticator{
		Username: c.config.User,
		Password: c.config.Pass,
	}

	if err := cluster.Authenticate(auth); err != nil {
		return err
	}

	c.conn, err = cluster.OpenBucket(c.config.Bucket, "")
	if err != nil {
		return err
	}

	return nil
}

// Close shutsdown a couchbase bucket connection
func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		c.conn = nil
		return err
	}

	c.conn = nil
	return nil
}

// Test opens a connection to a couchbase bucket and then closes it, returning an error
func (c *Client) Test() error {
	if c.conn != nil {
		c.Close()
	}

	if err := c.Connect(); err != nil {
		return err
	}
	defer c.Close()

	return nil
}

// EOF
