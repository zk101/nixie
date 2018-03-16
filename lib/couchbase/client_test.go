package couchbase

import "testing"

// TestNewClient checks NewClient
func TestNewClient(t *testing.T) {
	client := NewClient(nil)

	if client.config == nil {
		t.Error("client.config is nil, expect this to contain a valid pointer to Default Config")
	}
	if client.conn != nil {
		t.Error("client.conn is not nil, expect this to be nil")
	}
	if client.config.Cluster != "couchbase://localhost" {
		t.Errorf("Expected for client.config.Cluster value %s but got value %s", "couchbase://localhost", client.config.Cluster)
	}
	if client.config.Bucket != "" {
		t.Errorf("Expected for client.config.Bucket is an empty string but got value %s", client.config.Bucket)
	}
	if client.config.Pass != "" {
		t.Errorf("Expected for client.config.Pass is an empty string but got value %s", client.config.Pass)
	}

	config := Config{
		Cluster: "couchbase://cbtest.example.com",
		Bucket:  "TestBucket",
		Pass:    "password",
	}

	client = NewClient(&config)
	if client.config == nil {
		t.Error("Client config is nil, expect this to contain a valid pointer to Default Config")
	}
	if client.conn != nil {
		t.Error("client.conn is not nil, expect this to be nil")
	}
	if client.config.Cluster != "couchbase://cbtest.example.com" {
		t.Errorf("Expected for client.config.Cluster value %s but got value %s", "couchbase://cbtest.example.com", client.config.Cluster)
	}
	if client.config.Bucket != "TestBucket" {
		t.Errorf("Expected for client.config.Bucket value %sbut got value %s", "TestBucket", client.config.Bucket)
	}
	if client.config.Pass != "password" {
		t.Errorf("Expected for client.config.Pass value %s but got value %s", "password", client.config.Pass)
	}
}

// EOF
