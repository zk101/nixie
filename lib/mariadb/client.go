package mariadb

import "database/sql"

// Client holds operational values
type Client struct {
	config *Config
	conn   *sql.DB
}

// NewClient returns a config Client struct pointer
func NewClient(conf *Config) *Client {
	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	return &Client{
		config: conf,
	}
}

// EOF
