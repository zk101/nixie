package config

import (
	"errors"
	"sync"
	"time"
)

// Client contains logic of goroutine reuse.
type Client struct {
	interval  uint32
	config    interface{}
	prefix    string
	seperator string
	kv        kvstore
	wg        *sync.WaitGroup
	run       bool
}

// InitClient creates new goroutine pool with given config
func InitClient(interval uint32, config interface{}, prefix, seperator string, kv kvstore) (*Client, error) {
	if interval < 1 {
		return nil, errors.New("interval needs to be greater than 1")
	}

	if seperator == "" {
		return nil, errors.New("need a seperator")
	}

	if config == nil {
		return nil, errors.New("require a config struct pointer")
	}

	if kv == nil {
		return nil, errors.New("require a kvstore")
	}

	client := Client{
		interval:  interval,
		config:    config,
		prefix:    prefix,
		seperator: seperator,
		kv:        kv,
		wg:        &sync.WaitGroup{},
		run:       true,
	}

	go client.worker()

	time.Sleep(time.Millisecond * 5)

	return &client, nil
}

// worker exists to provide dynamic config overrides
func (c *Client) worker() {
	for c.run == true {
		Override(c.config, c.prefix, c.seperator, c.kv)

		time.Sleep(time.Millisecond * time.Duration(c.interval))
	}
}

// Stop shuts down a running client
func (c *Client) Stop() {
	c.run = false
	c.wg.Wait()
}

// EOF
