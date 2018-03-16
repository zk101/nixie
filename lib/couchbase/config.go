package couchbase

// Config contains couchbase connection details
type Config struct {
	Cluster string
	Bucket  string
	User    string
	Pass    string
}

// DefaultConfig sets up a Config struct
func DefaultConfig() Config {
	return Config{
		Cluster: "couchbase://localhost",
	}
}

// EOF
