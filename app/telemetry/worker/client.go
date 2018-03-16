package worker

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/telemetry/prometheus"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/mariadb"
	"github.com/zk101/nixie/lib/rabbitmq"
)

// Client holds worker config
type Client struct {
	config     *Config
	run        bool
	register   chan struct{}
	wg         sync.WaitGroup
	log        *logging.Client
	mqConfig   *rabbitmq.Config
	prometheus *prometheus.Config
	sql        *mariadb.Client
}

// NewClient sets the test up
func NewClient(conf *Config, mqConf *rabbitmq.Config, log *logging.Client, prometheus *prometheus.Config, sql *mariadb.Client) (*Client, error) {
	if mqConf == nil {
		return nil, errors.New("require a valid rabbitmq config")
	}

	if log == nil {
		return nil, errors.New("require a valid logging client")
	}

	if prometheus == nil {
		return nil, errors.New("require a valid prometheus config")
	}

	if sql == nil {
		return nil, errors.New("require a valid sql client")
	}

	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	if conf.PoolSize < 1 {
		return nil, errors.New("pool size must be 1 or greater")
	}

	client := Client{
		config:     conf,
		run:        true,
		register:   make(chan struct{}, conf.PoolSize),
		wg:         sync.WaitGroup{},
		log:        log,
		mqConfig:   mqConf,
		prometheus: prometheus,
		sql:        sql,
	}

	go client.monitor()

	time.Sleep(time.Millisecond * 5)

	return &client, nil
}

// Stop shuts down a running client
func (c *Client) Stop() {
	c.run = false
	c.wg.Wait()
}

// monitor ensures min workers are running
func (c *Client) monitor() {
	for c.run == true {
		for x := len(c.register); x < c.config.PoolSize; x++ {
			c.register <- struct{}{}
			go c.worker(uuid.NewV4().String())
		}

		time.Sleep(time.Millisecond * 50)
	}
}

// EOF
