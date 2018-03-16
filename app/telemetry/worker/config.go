package worker

// Config holds worker pool config options
type Config struct {
	PoolSize     int
	WorkerExpiry int
}

// DefaultConfig returns a Config struct with predefined settings
func DefaultConfig() Config {
	return Config{
		PoolSize:     10,
		WorkerExpiry: 30,
	}
}

// EOF
