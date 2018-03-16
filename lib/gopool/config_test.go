package gopool

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.PoolSizeMin != 1 {
		t.Errorf("Expected for config.Min value %d but got value %d", 1, config.PoolSizeMin)
	}
	if config.PoolSizeMax != 10 {
		t.Errorf("Expected for config.Max value %d but got value %d", 10, config.PoolSizeMax)
	}
	if config.QueueSize != 10 {
		t.Errorf("Expected for config.QueueSize value %d but got value %d", 10, config.QueueSize)
	}
	if config.ScheduleTimeout != 50 {
		t.Errorf("Expected for config.Timeout value %d but got value %d", 50, config.ScheduleTimeout)
	}
	if config.WorkerExpiry != 30 {
		t.Errorf("Expected for config.Expiry value %d but got value %d", 30, config.WorkerExpiry)
	}
}

// EOF
