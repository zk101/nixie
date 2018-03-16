package httpd

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zk101/nixie/lib/couchbase"
	"github.com/zk101/nixie/lib/ldap"
	"github.com/zk101/nixie/lib/network"
	"github.com/zk101/nixie/lib/rabbitmq"
)

// Start sets up and starts the httpd server
func (c *Client) Start() error {
	http.HandleFunc("/", http.NotFound)
	http.HandleFunc("/health", c.handlerHealth)
	http.Handle("/metrics", promhttp.Handler())

	c.server = &http.Server{
		Addr: ":" + strconv.Itoa(c.config.Port),
	}
	c.server.SetKeepAlivesEnabled(c.config.Keepalive)

	if c.config.SSL == true {
		if c.config.SSLkey == "" {
			return errors.New("missing ssl key")
		}

		if _, err := os.Stat(c.config.SSLkey); err != nil {
			return err
		}

		if c.config.SSLcert == "" {
			return errors.New("missing ssl cert")
		}

		if _, err := os.Stat(c.config.SSLcert); err != nil {
			return err
		}
	}

	go func() {
		if c.config.SSL == true {
			c.server.ListenAndServeTLS(c.config.SSLcert, c.config.SSLkey)
		} else {
			c.server.ListenAndServe()
		}
	}()

	return nil
}

// Stop sets up and starts the httpd server
func (c *Client) Stop() error {
	return c.server.Close()
}

// handlerHealth deals with health requests
func (c *Client) handlerHealth(w http.ResponseWriter, r *http.Request) {
	if err := network.CheckHTTPhost(r.RemoteAddr, `^((127\.0\.0\.)|(10\.)|(192\.168\.)|(172\.1[6-9]\.)|(172\.2[0-9]\.)|(172\.3[0-1]\.))`); err != nil {
		c.log.Sugar().Errorw("remote address test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	testCB := couchbase.NewClient(c.cbPool.GetBucketConfig())
	if err := testCB.Test(); err != nil {
		c.log.Sugar().Errorw("couchbase test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		http.Error(w, "503 Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	testLDAPro := ldap.NewClient(c.ldapPool.GetROconfig(), c.ldapPool.GetCAcertPool())
	if err := testLDAPro.Test(); err != nil {
		c.log.Sugar().Errorw("ldap readonly test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		http.Error(w, "503 Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	testLDAPrw := ldap.NewClient(c.ldapPool.GetRWconfig(), c.ldapPool.GetCAcertPool())
	if err := testLDAPrw.Test(); err != nil {
		c.log.Sugar().Errorw("ldap readwrite test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		http.Error(w, "503 Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	testRabbitMQ, err := rabbitmq.NewClient(c.mqConfig)
	if err != nil {
		c.log.Sugar().Errorw("rabbitmq test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		http.Error(w, "503 Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	defer testRabbitMQ.Close()

	http.Error(w, "200 Okay", http.StatusOK)
}

// EOF
