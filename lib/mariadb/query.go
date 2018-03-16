package mariadb

import (
	"database/sql"
)

// Query runs an sql query
func (c *Client) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	rows, err := c.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// EOF
