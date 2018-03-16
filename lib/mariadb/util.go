package mariadb

import (
	"database/sql"

	// mysql is the driver used to open a connection
	_ "github.com/go-sql-driver/mysql"
)

// Connect starts an sql connection
func (c *Client) Connect() error {
	if c.conn != nil {
		c.Close()
	}

	var err error
	c.conn, err = sql.Open("mysql", c.config.DSN)
	if err != nil {
		return err
	}
	c.conn.SetMaxIdleConns(c.config.MaxIdle)
	c.conn.SetMaxOpenConns(c.config.MaxOpen)

	return nil
}

// Close stops an sql connection
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

	return c.conn.Ping()
}

// EOF
