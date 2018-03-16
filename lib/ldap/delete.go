package ldap

import "github.com/mavricknz/ldap"

// Delete removes an object from ldap
func (c *Client) Delete(delReq *ldap.DeleteRequest) error {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	if err := c.conn.Delete(delReq); err != nil {
		return err
	}

	return nil
}

// EOF
