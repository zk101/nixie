package rabbitmq

// Config holds configuration data
type Config struct {
	URL string
}

// DefaultConfig returns a Config struct with default values
func DefaultConfig() Config {
	return Config{
		URL: "amqp://guest:guest@tcp(127.0.0.1):5672/",
	}
}

// EOF
