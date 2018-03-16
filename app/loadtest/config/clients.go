package config

import (
	"crypto/x509"
	"errors"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/httpd"
	"github.com/zk101/nixie/app/loadtest/prometheus"
	"github.com/zk101/nixie/app/loadtest/websocket"
	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/consul"
	"github.com/zk101/nixie/lib/logging"
)

// Clients holds pointers to clients used in the application
type Clients struct {
	ServiceID  string
	CAcert     *x509.CertPool
	Consul     *consul.Client
	HTTPD      *http.Server
	Log        *logging.Client
	Prometheus *prometheus.Config
	WebSocket  *websocket.Client
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

	clients.HTTPD, err = httpd.StartHTTPD(&conf.HTTPD)
	if err != nil {
		return nil, err
	}

	return &clients, nil
}

// ShutdownClients stops all the Clients
func ShutdownClients(clients *Clients) {
	httpd.StopHTTPD(clients.HTTPD)
	clients.Log.Sync()
}

// EOF
