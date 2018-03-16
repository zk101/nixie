package mariadb

// Config holds configuration values
type Config struct {
	DSN     string
	MaxIdle int
	MaxOpen int
}

// DefaultConfig returns a Config struct with default values
func DefaultConfig() Config {
	return Config{
		DSN:     "root:password@localhost/test",
		MaxIdle: 1,
		MaxOpen: 10,
	}
}

// EOF
