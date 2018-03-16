package httpd

import (
	"errors"
	"net/http"

	"github.com/zk101/nixie/lib/couchbase"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/rabbitmq"
)

// Client adds operation struct to http which supports additional pointers
type Client struct {
	serviceID  string
	config     *httpd.Config
	cbPresence *couchbase.Config
	mqConfig   *rabbitmq.Config
	log        *logging.Client
	server     *http.Server
}

// NewClient returns a configured Client pointer
func NewClient(serviceID string, conf *httpd.Config, cbPresence *couchbase.Config, mqConfig *rabbitmq.Config, log *logging.Client) (*Client, error) {
	if serviceID == "" {
		return nil, errors.New("require a none empty service id")
	}

	if conf == nil {
		return nil, errors.New("require an http config")
	}

	if cbPresence == nil {
		return nil, errors.New("require a couchbase config")
	}

	if mqConfig == nil {
		return nil, errors.New("require a rabbitmq config")
	}

	if log == nil {
		return nil, errors.New("require a logging client")
	}

	client := Client{
		serviceID:  serviceID,
		config:     conf,
		cbPresence: cbPresence,
		mqConfig:   mqConfig,
		log:        log,
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return &client, nil
}

// EOF
