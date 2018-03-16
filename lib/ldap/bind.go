package ldap

import (
	"crypto/tls"

	"github.com/mavricknz/ldap"
)

// Bind tests if a given set of credentials can bind on an ldap connection
func (c *Client) Bind(username, password string) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	return c.conn.Bind(username, password)
}

// CheckBind tests if a given set of credentials can bind to an ldap server
func (c *Client) CheckBind(user, pass string) error {
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
	defer c.Close()

	return c.conn.Bind(user, pass)
}

// EOF
