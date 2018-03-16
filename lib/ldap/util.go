package ldap

import (
	"crypto/tls"

	"github.com/mavricknz/ldap"
)

// Connect opens a connection to an LDAP server defined in Config
func (c *Client) Connect() error {
	if c.conn != nil {
		if err := c.Close(); err != nil {
			return err
		}
	}

	if c.config.SSL == true {
		c.conn = ldap.NewLDAPSSLConnection(c.config.Host, c.config.Port, &tls.Config{ServerName: c.config.Host, RootCAs: c.cacert})
	} else {
		c.conn = ldap.NewLDAPConnection(c.config.Host, c.config.Port)
	}

	if err := c.conn.Connect(); err != nil {
		return err
	}

	if err := c.conn.Bind(c.config.User, c.config.Pass); err != nil {
		c.Close()
		return err
	}

	return nil
}

// Close shuts down an LDAP connection
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

// Test checks a connection can be formed to an ldap server
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
