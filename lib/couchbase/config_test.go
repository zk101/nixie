package couchbase

import "testing"

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Cluster != "couchbase://localhost" {
		t.Errorf("Expected for config.Cluster value %s but got value %s", "couchbase://localhost", config.Cluster)
	}
	if config.Bucket != "" {
		t.Errorf("Expected for config.Bucket is an empty string but got value %s", config.Bucket)
	}
	if config.User != "" {
		t.Errorf("Expected for config.User is an empty string but got value %s", config.User)
	}
	if config.Pass != "" {
		t.Errorf("Expected for config.Pass is an empty string but got value %s", config.Pass)
	}
}

// EOF
