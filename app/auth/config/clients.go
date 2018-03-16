package config

import (
	"crypto/x509"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/auth/prometheus"
	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/consul"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	ldappool "github.com/zk101/nixie/lib/ldap/pool"
	"github.com/zk101/nixie/lib/logging"
)

// Clients holds pointers to clients used in the application
type Clients struct {
	ServiceID  string
	CAcert     *x509.CertPool
	CBpool     *cbpool.Client
	Consul     *consul.Client
	LDAPpool   *ldappool.Client
	Log        *logging.Client
	Prometheus *prometheus.Config
}

// SetupClients starts all the clients
func SetupClients(conf *Config) (*Clients, error) {
	if conf == nil {
		return nil, errors.New("setupclients requires a valid conf")
	}

	var err error
	clients := Clients{
		ServiceID: uuid.NewV4().String(),
	}

	clients.Log, err = logging.NewClient(&conf.Log, clients.ServiceID)
	if err != nil {
		return nil, err
	}

	clients.CAcert, err = config.LoadCAcerts(conf.Controls.CAcertPath)
	if err != nil {
		return nil, err
	}

	clients.Consul = consul.NewClient(&conf.Consul)

	clients.Prometheus, err = prometheus.Init()
	if err != nil {
		return nil, err
	}

	clients.LDAPpool, err = ldappool.InitClient(&conf.LDAPpool, &conf.LDAPro, &conf.LDAPrw, clients.CAcert)
	if err != nil {
		return nil, err
	}

	clients.CBpool, err = cbpool.InitClient(&conf.CBpool, &conf.CBpresence)
	if err != nil {
		clients.LDAPpool.Stop()
		return nil, err
	}

	return &clients, nil
}

// ShutdownClients stops all the Clients
func ShutdownClients(clients *Clients) {
	clients.CBpool.Stop()
	clients.LDAPpool.Stop()
	clients.Log.Sync()
}

// EOF
