package gelf

import "testing"

func TestDefaultConfig(t *testing.T) {
	result := DefaultConfig("foo")

	if result.Host != "foo" {
		t.Errorf("Expected for result.Host value %v but got value %v", "foo", result.Host)
	}
	if result.Port != 12201 {
		t.Errorf("Expected for result.Port value %v but got value %v", 8154, result.Port)
	}
	if result.MaxChunkSize != 8154 {
		t.Errorf("Expected for result.MaxChunkSize value %v but got value %v", 8154, result.MaxChunkSize)
	}
	if result.Compression != CompressionNone {
		t.Errorf("Expected for result.Compression value %v but got value %v", CompressionNone, result.Compression)
	}
}

// EOF
