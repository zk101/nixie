package config

import (
	"github.com/zk101/nixie/app/chat/worker"
	"github.com/zk101/nixie/lib/consul"
	"github.com/zk101/nixie/lib/couchbase"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/ldap"
	ldappool "github.com/zk101/nixie/lib/ldap/pool"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/rabbitmq"
	"github.com/zk101/nixie/lib/signal"
)

// Config holds core configuration
type Config struct {
	CBpool     cbpool.Config
	CBpresence couchbase.Config
	Consul     consul.Config
	Controls   Controls
	HTTPD      httpd.Config
	LDAPpool   ldappool.Config
	LDAPro     ldap.Config
	LDAPrw     ldap.Config
	Log        logging.Config
	RabbitMQ   rabbitmq.Config
	Signal     signal.Config
	Worker     worker.Config
}

// NewConfig creates a Config struct and returns a pointer
func NewConfig() *Config {
	return &Config{
		CBpool:     cbpool.DefaultConfig(),
		CBpresence: couchbase.DefaultConfig(),
		Consul:     consul.DefaultConfig(),
		Controls:   DefaultControls(),
		HTTPD:      httpd.DefaultConfig(),
		LDAPpool:   ldappool.DefaultConfig(),
		LDAPro:     ldap.DefaultConfig(),
		LDAPrw:     ldap.DefaultConfig(),
		Log:        logging.DefaultConfig(),
		RabbitMQ:   rabbitmq.DefaultConfig(),
		Signal:     signal.DefaultConfig(),
		Worker:     worker.DefaultConfig(),
	}
}

// EOF
