package ldap

import "github.com/mavricknz/ldap"

// Add inserts a new object to ldap
func (c *Client) Add(req *ldap.AddRequest) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	if err := c.conn.Add(req); err != nil {
		return err
	}

	return nil
}

// EOF
