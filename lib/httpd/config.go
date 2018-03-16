package httpd

// Config holds httpd details
type Config struct {
	Port      int
	Keepalive bool
	SSL       bool
	SSLkey    string
	SSLcert   string
}

// DefaultConfig sets up a new Config instance
func DefaultConfig() Config {
	return Config{
		Port: 443,
	}
}

// EOF
