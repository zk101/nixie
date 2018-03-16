package ldap

// Config struct holds LDAP connection details
type Config struct {
	Host string
	Port uint16
	User string
	Pass string
	Base string
	SSL  bool
}

// DefaultConfig sets up an Ldap Config struct
func DefaultConfig() Config {
	return Config{
		Host: "localhost",
		Port: 389,
		Base: "dc=example,dc=org",
	}
}

// EOF
