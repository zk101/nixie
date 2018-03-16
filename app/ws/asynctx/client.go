package asynctx

import (
	"errors"
	"sync"
	"time"

	"github.com/zk101/nixie/app/ws/prometheus"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/rabbitmq"
)

// Client holds operation values for the worker pools
type Client struct {
	serviceID  string
	config     *Config
	mqConfig   *rabbitmq.Config
	register   chan struct{}
	work       chan task
	wg         *sync.WaitGroup
	run        bool
	log        *logging.Client
	prometheus *prometheus.Config
}

// NewClient sets up new worker pools
func NewClient(serviceID string, conf *Config, mqConfig *rabbitmq.Config, log *logging.Client, prometheus *prometheus.Config) (*Client, error) {
	if serviceID == "" {
		return nil, errors.New("service id can not be empty")
	}

	if mqConfig == nil {
		return nil, errors.New("missing rabbitmq config")
	}

	if log == nil {
		return nil, errors.New("missing logging client")
	}

	if prometheus == nil {
		return nil, errors.New("missing prometheus client")
	}

	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	if conf.PoolSizeMin < 1 {
		return nil, errors.New("pool spawn min must be 1 or greater")
	}

	if conf.PoolSizeMin > conf.PoolSizeMax {
		return nil, errors.New("pool spawn min must not be greater than pool spawn max")
	}

	client := Client{
		serviceID:  serviceID,
		config:     conf,
		mqConfig:   mqConfig,
		register:   make(chan struct{}, conf.PoolSizeMax),
		work:       make(chan task, 10),
		wg:         &sync.WaitGroup{},
		run:        true,
		log:        log,
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
		for x := len(c.register); x < c.config.PoolSizeMin; x++ {
			c.register <- struct{}{}
			go c.worker(task{}, true)
		}

		c.prometheus.SetAsynctxQueueCount(len(c.work))
		c.prometheus.SetAsynctxWorkerCount(len(c.register))

		time.Sleep(time.Millisecond * 50)
	}
}

// EOF
