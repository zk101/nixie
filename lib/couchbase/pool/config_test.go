package pool

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Min != 1 {
		t.Errorf("Expected for config.Min value %d but got value %d", 1, config.Min)
	}
	if config.Max != 10 {
		t.Errorf("Expected for config.Max value %d but got value %d", 10, config.Max)
	}
	if config.QueueSize != 10 {
		t.Errorf("Expected for config.QueueSize value %d but got value %d", 10, config.QueueSize)
	}
	if config.Timeout != 50 {
		t.Errorf("Expected for config.Timeout value %d but got value %d", 50, config.Timeout)
	}
	if config.Expiry != 30 {
		t.Errorf("Expected for config.Expiry value %d but got value %d", 30, config.Expiry)
	}
}

// EOF
