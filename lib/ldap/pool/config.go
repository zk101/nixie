package pool

// Config holds static configuration options
type Config struct {
	Min       int
	Max       int
	QueueSize int
	Timeout   int
	Expiry    int
}

// DefaultConfig returns a struct with default Config values
func DefaultConfig() Config {
	return Config{
		Min:       1,
		Max:       10,
		QueueSize: 10,
		Timeout:   50,
		Expiry:    30,
	}
}

// EOF
