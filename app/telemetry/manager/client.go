package manager

import (
	"errors"
	"sync"
	"time"

	"github.com/zk101/nixie/app/telemetry/prometheus"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/mariadb"
)

// Client holds worker config
type Client struct {
	run        bool
	wg         sync.WaitGroup
	log        *logging.Client
	prometheus *prometheus.Config
	sql        *mariadb.Client
}

// NewClient sets the test up
func NewClient(log *logging.Client, prometheus *prometheus.Config, sql *mariadb.Client) (*Client, error) {
	if log == nil {
		return nil, errors.New("require a valid logging client")
	}

	if prometheus == nil {
		return nil, errors.New("require a valid prometheus config")
	}

	if prometheus == nil {
		return nil, errors.New("require a valid prometheus config")
	}

	if sql == nil {
		return nil, errors.New("require a valid sql client")
	}

	client := Client{
		run:        true,
		wg:         sync.WaitGroup{},
		log:        log,
		prometheus: prometheus,
		sql:        sql,
	}

	go client.manager()
	time.Sleep(time.Millisecond * 5)

	return &client, nil
}

// Stop shuts down a running client
func (c *Client) Stop() {
	c.run = false
	c.wg.Wait()
}

// EOF
