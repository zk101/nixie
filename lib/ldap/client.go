package ldap

import (
	"crypto/x509"

	"github.com/mavricknz/ldap"
)

// Client is an operational struct which all ldap methods are based
type Client struct {
	config *Config
	cacert *x509.CertPool
	conn   *ldap.LDAPConnection
}

// NewClient returns a pointer to a new Client instance
func NewClient(conf *Config, cacert *x509.CertPool) *Client {
	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	return &Client{
		config: conf,
		cacert: cacert,
	}
}

// EOF
