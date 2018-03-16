package worker

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/chat/prometheus"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	ldappool "github.com/zk101/nixie/lib/ldap/pool"
	"github.com/zk101/nixie/lib/logging"
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
	cbPool     *cbpool.Client
	ldapPool   *ldappool.Client
	prometheus *prometheus.Config
	mqClient   *rabbitmq.Client
}

// NewClient sets the test up
func NewClient(conf *Config, mqConf *rabbitmq.Config, cbPool *cbpool.Client, ldapPool *ldappool.Client, log *logging.Client, prometheus *prometheus.Config) (*Client, error) {
	if mqConf == nil {
		return nil, errors.New("require a valid rabbitmq config")
	}

	if cbPool == nil {
		return nil, errors.New("require a valid couchbase pool client")
	}

	if ldapPool == nil {
		return nil, errors.New("require a valid ldap pool client")
	}

	if log == nil {
		return nil, errors.New("require a valid logging client")
	}

	if prometheus == nil {
		return nil, errors.New("require a valid prometheus config")
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
		cbPool:     cbPool,
		ldapPool:   ldapPool,
		prometheus: prometheus,
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
