package ldap

import "github.com/mavricknz/ldap"

// Modify updates an object in ldap
func (c *Client) Modify(modReq *ldap.ModifyRequest) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	if err := c.conn.Modify(modReq); err != nil {
		return err
	}

	return nil
}

// ModDn updates an objects DN in ldap
func (c *Client) ModDn(req *ldap.ModDnRequest) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	if err := c.conn.ModDn(req); err != nil {
		return err
	}

	return nil
}

// EOF
