package pool

import (
	"crypto/x509"
	"errors"
	"sync"
	"time"

	"github.com/zk101/nixie/lib/ldap"
)

// Client holds operational values for the pool
type Client struct {
	config   *Config
	cacert   *x509.CertPool
	ldapbind *ldap.Config
	ldapro   *ldap.Config
	ldaprw   *ldap.Config
	register chan struct{}
	work     chan taskRequest
	wg       *sync.WaitGroup
	run      bool
}

// InitClient sets up a new client and returns a pointer to it
func InitClient(conf *Config, ldapro *ldap.Config, ldaprw *ldap.Config, cacert *x509.CertPool) (*Client, error) {
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

	if ldapro == nil && ldaprw == nil {
		return nil, errors.New("no ldap config given require ldapro or ldaprw or both")
	}

	if ldapro == nil {
		ldapro = ldaprw
	}

	if ldapro.Base != ldaprw.Base {
		return nil, errors.New("base for ldapro and ldaprw must match")
	}

	client := Client{
		config:   conf,
		cacert:   cacert,
		ldapbind: ldapro,
		ldapro:   ldapro,
		ldaprw:   ldaprw,
		register: make(chan struct{}, conf.Max),
		work:     make(chan taskRequest, conf.QueueSize),
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

// GetBase returns the Base from LDAP config
func (c *Client) GetBase() string {
	return c.ldapro.Base
}

// GetROconfig returns ldap.Config pointer for RO configuration
func (c *Client) GetROconfig() *ldap.Config {
	return c.ldapro
}

// GetRWconfig returns ldap.Config pointer for RW configuration
func (c *Client) GetRWconfig() *ldap.Config {
	return c.ldaprw
}

// GetCAcertPool returns x509 CertPool used in the pool
func (c *Client) GetCAcertPool() *x509.CertPool {
	return c.cacert
}

// EOF
