package logging

import "testing"

func TestDefaultConfig(t *testing.T) {
	result := DefaultConfig()

	if result.Method != "console" {
		t.Errorf("Expected for result.Method value %v but got value %v", "console", result.Method)
	}
}

// EOF
