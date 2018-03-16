package consul

// Config holds consul details
type Config struct {
	Address    string
	Scheme     string
	Datacenter string
	Token      string
}

// DefaultConfig setups a new Config struct and returns it
func DefaultConfig() Config {
	return Config{
		Address:    "localhost:8500",
		Scheme:     "http",
		Datacenter: "global",
		Token:      "",
	}
}

// EOF
