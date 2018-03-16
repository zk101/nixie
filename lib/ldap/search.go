package ldap

import (
	"github.com/mavricknz/ldap"
)

// Search runs a filter on ldap and returns the search handle
func (c *Client) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	sr, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

// EOF
