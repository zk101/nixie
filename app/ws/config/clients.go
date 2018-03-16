package config

import (
	"crypto/x509"
	"errors"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/ws/asyncrx"
	"github.com/zk101/nixie/app/ws/asynctx"
	"github.com/zk101/nixie/app/ws/connection"
	"github.com/zk101/nixie/app/ws/httpd"
	"github.com/zk101/nixie/app/ws/prometheus"
	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/consul"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/gopool"
	"github.com/zk101/nixie/lib/logging"
)

// Clients holds pointers to clients used in the application
type Clients struct {
	ServiceID  string
	AsyncRX    *asyncrx.Client
	AsyncTX    *asynctx.Client
	CAcert     *x509.CertPool
	CBpool     *cbpool.Client
	Connection *connection.Client
	Consul     *consul.Client
	GoPool     *gopool.Client
	HTTPD      *httpd.Client
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

	clients.AsyncTX, err = asynctx.NewClient(clients.ServiceID, &conf.AsyncTX, &conf.RabbitMQ, clients.Log, clients.Prometheus)
	if err != nil {
		return nil, err
	}

	clients.HTTPD, err = httpd.NewClient(clients.ServiceID, &conf.HTTPD, &conf.CBpresence, &conf.RabbitMQ, clients.Log)
	if err != nil {
		return nil, err
	}

	clients.CBpool, err = cbpool.InitClient(&conf.CBpool, &conf.CBpresence)
	if err != nil {
		clients.HTTPD.Stop()
		return nil, err
	}

	clients.GoPool, err = gopool.InitClient(&conf.GoPool)
	if err != nil {
		clients.HTTPD.Stop()
		clients.CBpool.Stop()
		return nil, err
	}

	return &clients, nil
}

// ShutdownClients stops all the Clients
func ShutdownClients(clients *Clients) {
	clients.AsyncRX.Stop()
	clients.HTTPD.Stop()
	clients.AsyncTX.Stop()
	clients.CBpool.Stop()
	clients.Log.Sync()
}

// EOF
