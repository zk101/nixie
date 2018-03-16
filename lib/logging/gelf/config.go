package gelf

// Config represents the required settings for connecting the gelf data sink.
type Config struct {
	Host         string
	Port         int
	MaxChunkSize int
	Compression  int
}

// DefaultConfig provides a configuration with default values for port and chunk size.
func DefaultConfig(host string) Config {
	return Config{
		Host:         host,
		Port:         12201,
		MaxChunkSize: 8154,
		Compression:  CompressionNone,
	}
}

// EOF
