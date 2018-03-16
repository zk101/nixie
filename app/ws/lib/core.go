package lib

import (
	"errors"
	"sync"
	"time"

	"github.com/zk101/nixie/app/ws/asyncrx"

	"github.com/zk101/nixie/app/ws/config"
	"github.com/zk101/nixie/app/ws/connection"
)

// Core holds pointers to Core Config objects
type Core struct {
	Config  *config.Config
	Clients *config.Clients
	run     bool
	wg      *sync.WaitGroup
}

// NewCore returns a setup client pointer or nil
func NewCore(conf *config.Config) (*Core, error) {
	if conf == nil {
		return nil, errors.New("newcore requires a valid config")
	}

	core := Core{
		Config: conf,
		run:    true,
		wg:     &sync.WaitGroup{},
	}

	var err error

	core.Clients, err = config.SetupClients(conf)
	if err != nil {
		return nil, err
	}

	if err := core.SetupQueues(); err != nil {
		return nil, err
	}

	core.Clients.Connection, err = connection.NewClient(core.Clients.ServiceID, core.Clients.Log, core.Clients.Prometheus, core.Clients.AsyncTX, core.Clients.CBpool)
	if err != nil {
		return nil, err
	}

	core.Clients.AsyncRX, err = asyncrx.NewClient(core.Clients.ServiceID, &conf.AsyncRX, &conf.RabbitMQ, core.Clients.Log, core.Clients.Prometheus, core.Clients.Connection.GetUserMap())
	if err != nil {
		return nil, err
	}

	if err := core.doConsulReg(); err != nil {
		return nil, err
	}

	go core.Clients.Prometheus.Worker(&core.Config.Signal, core.Clients.CBpool, core.Clients.GoPool)

	go core.worker()
	time.Sleep(time.Millisecond * 5)

	return &core, nil
}

// Stop shuts down the application
func (c *Core) Stop() error {
	c.run = false

	return c.doConsulUnreg()
}

// EOF
