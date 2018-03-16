package config

import (
	"github.com/zk101/nixie/app/ws/asyncrx"
	"github.com/zk101/nixie/app/ws/asynctx"
	"github.com/zk101/nixie/lib/consul"
	"github.com/zk101/nixie/lib/couchbase"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/gopool"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/rabbitmq"
	"github.com/zk101/nixie/lib/signal"
)

// Config holds core configuration
type Config struct {
	AsyncRX    asyncrx.Config
	AsyncTX    asynctx.Config
	CBpool     cbpool.Config
	CBpresence couchbase.Config
	Consul     consul.Config
	Controls   Controls
	GoPool     gopool.Config
	HTTPD      httpd.Config
	Log        logging.Config
	RabbitMQ   rabbitmq.Config
	Signal     signal.Config
}

// NewConfig creates a Config struct and returns a pointer
func NewConfig() *Config {
	return &Config{
		AsyncRX:    asyncrx.DefaultConfig(),
		AsyncTX:    asynctx.DefaultConfig(),
		CBpool:     cbpool.DefaultConfig(),
		CBpresence: couchbase.DefaultConfig(),
		Consul:     consul.DefaultConfig(),
		Controls:   DefaultControls(),
		GoPool:     gopool.DefaultConfig(),
		HTTPD:      httpd.DefaultConfig(),
		Log:        logging.DefaultConfig(),
		RabbitMQ:   rabbitmq.DefaultConfig(),
		Signal:     signal.DefaultConfig(),
	}
}

// EOF
