package pool

import (
	"errors"
	"sync"
	"time"

	"github.com/zk101/nixie/lib/couchbase"
)

// Client holds operational values for the pool
type Client struct {
	config   *Config
	bucket   *couchbase.Config
	register chan struct{}
	work     chan taskRequest
	wg       *sync.WaitGroup
	run      bool
}

// InitClient sets up a new client and returns a pointer to it
func InitClient(conf *Config, bucket *couchbase.Config) (*Client, error) {
	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	if conf.Min < 1 {
		return nil, errors.New("pool spawn min must be 1 or greater")
	}

	if conf.Min > conf.Max {
		return nil, errors.New("pool spawn min must not be greater than pool spawn max")
	}

	if bucket == nil {
		return nil, errors.New("no bucket config given")
	}

	client := Client{
		config:   conf,
		bucket:   bucket,
		register: make(chan struct{}, conf.Max),
		work:     make(chan taskRequest, conf.QueueSize),
		wg:       &sync.WaitGroup{},
		run:      true,
	}

	clientCBtest := couchbase.NewClient(bucket)
	if err := clientCBtest.Test(); err != nil {
		return nil, err
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

// scheduleWrapper provides an entry point to queuing new tasks
func (c *Client) scheduleWrapper(task taskRequest) error {
	return c.schedule(task, time.After(time.Millisecond*time.Duration(c.config.Timeout)))
}

// schedule does the actual scheduling
func (c *Client) schedule(task taskRequest, timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return errors.New("timeout scheduling task")
	case c.work <- task:
		return nil
	case c.register <- struct{}{}:
		go c.worker(task, false)
		return nil
	}
}

// monitor ensures min workers are running
func (c *Client) monitor() {
	for c.run == true {
		for x := len(c.register); x < c.config.Min; x++ {
			c.register <- struct{}{}
			go c.worker(taskRequest{}, true)
		}

		time.Sleep(time.Millisecond * 50)
	}
}

// GetWorkerCount returns the current number of workers
func (c *Client) GetWorkerCount() int {
	return len(c.register)
}

// GetQueueCount returns the current number of tasks on the queue
func (c *Client) GetQueueCount() int {
	return len(c.work)
}

// GetBucketConfig returns a couchbase.Config pointer to the Bucket config
func (c *Client) GetBucketConfig() *couchbase.Config {
	return c.bucket
}

// EOF
