package gopool

// Config holds async worker pool config options
type Config struct {
	PoolSizeMin     int
	PoolSizeMax     int
	QueueSize       int
	ScheduleTimeout int
	WorkerExpiry    int
}

// DefaultConfig returns a Config struct with predefined settings
func DefaultConfig() Config {
	return Config{
		PoolSizeMin:     1,
		PoolSizeMax:     10,
		QueueSize:       10,
		ScheduleTimeout: 50,
		WorkerExpiry:    30,
	}
}

// EOF
