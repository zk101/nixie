package gopool

import (
	"errors"
	"sync"
	"time"
)

// Client contains logic of goroutine reuse.
type Client struct {
	config   *Config
	register chan struct{}
	work     chan func()
	wg       *sync.WaitGroup
	run      bool
}

// InitClient creates new goroutine pool with given config
func InitClient(conf *Config) (*Client, error) {
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
		config:   conf,
		register: make(chan struct{}, conf.PoolSizeMax),
		work:     make(chan func(), conf.QueueSize),
		wg:       &sync.WaitGroup{},
		run:      true,
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
			go c.worker(nil)
		}

		time.Sleep(time.Millisecond * 50)
	}
}

// EOF
