package mariadb

import (
	"database/sql"
)

// Exec wraps the db.Exec function
func (c *Client) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}

	result, err := c.conn.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// EOF
