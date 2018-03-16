package config

import (
	"github.com/zk101/nixie/app/telemetry/worker"
	"github.com/zk101/nixie/lib/consul"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/mariadb"
	"github.com/zk101/nixie/lib/rabbitmq"
	"github.com/zk101/nixie/lib/signal"
)

// Config holds core configuration
type Config struct {
	Consul     consul.Config
	Controls   Controls
	HTTPD      httpd.Config
	Log        logging.Config
	RabbitMQ   rabbitmq.Config
	Signal     signal.Config
	SQLmanager mariadb.Config
	SQLworker  mariadb.Config
	Worker     worker.Config
}

// NewConfig creates a Config struct and returns a pointer
func NewConfig() *Config {
	return &Config{
		Consul:     consul.DefaultConfig(),
		Controls:   DefaultControls(),
		HTTPD:      httpd.DefaultConfig(),
		Log:        logging.DefaultConfig(),
		RabbitMQ:   rabbitmq.DefaultConfig(),
		Signal:     signal.DefaultConfig(),
		SQLmanager: mariadb.DefaultConfig(),
		SQLworker:  mariadb.DefaultConfig(),
		Worker:     worker.DefaultConfig(),
	}
}

// EOF
