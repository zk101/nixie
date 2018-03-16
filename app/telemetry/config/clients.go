package config

import (
	"crypto/x509"
	"errors"

	"github.com/zk101/nixie/lib/mariadb"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/telemetry/httpd"
	"github.com/zk101/nixie/app/telemetry/manager"
	"github.com/zk101/nixie/app/telemetry/prometheus"
	"github.com/zk101/nixie/app/telemetry/worker"
	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/consul"
	"github.com/zk101/nixie/lib/logging"
)

// Clients holds pointers to clients used in the application
type Clients struct {
	ServiceID  string
	CAcert     *x509.CertPool
	Consul     *consul.Client
	HTTPD      *httpd.Client
	Log        *logging.Client
	Manager    *manager.Client
	Prometheus *prometheus.Config
	SQLmanager *mariadb.Client
	SQLworker  *mariadb.Client
	Worker     *worker.Client
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

	clients.SQLmanager = mariadb.NewClient(&conf.SQLmanager)
	if err := clients.SQLmanager.Connect(); err != nil {
		return nil, err
	}

	clients.SQLworker = mariadb.NewClient(&conf.SQLworker)
	if err := clients.SQLworker.Connect(); err != nil {
		return nil, err
	}

	clients.HTTPD, err = httpd.NewClient(&conf.HTTPD, &conf.RabbitMQ, &conf.SQLmanager, &conf.SQLworker, clients.Log)
	if err != nil {
		return nil, err
	}

	clients.Manager, err = manager.NewClient(clients.Log, clients.Prometheus, clients.SQLmanager)
	if err != nil {
		return nil, err
	}

	clients.Worker, err = worker.NewClient(&conf.Worker, &conf.RabbitMQ, clients.Log, clients.Prometheus, clients.SQLworker)
	if err != nil {
		clients.Manager.Stop()
		return nil, err
	}

	return &clients, nil
}

// ShutdownClients stops all the Clients
func ShutdownClients(clients *Clients) {
	clients.HTTPD.Stop()
	clients.Worker.Stop()
	clients.Manager.Stop()
	clients.Log.Sync()
}

// EOF
