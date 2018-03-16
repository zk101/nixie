package httpd

import (
	"errors"
	"net/http"

	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/httpd"
	ldappool "github.com/zk101/nixie/lib/ldap/pool"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/rabbitmq"
)

// Client adds operation struct to http which supports additional pointers
type Client struct {
	config   *httpd.Config
	cbPool   *cbpool.Client
	ldapPool *ldappool.Client
	mqConfig *rabbitmq.Config
	log      *logging.Client
	server   *http.Server
}

// NewClient returns a configured Client pointer
func NewClient(conf *httpd.Config, cbPool *cbpool.Client, ldapPool *ldappool.Client, mqConfig *rabbitmq.Config, log *logging.Client) (*Client, error) {
	if conf == nil {
		return nil, errors.New("require an http config")
	}

	if cbPool == nil {
		return nil, errors.New("require couchbase pool")
	}

	if ldapPool == nil {
		return nil, errors.New("require ldap pool")
	}

	if mqConfig == nil {
		return nil, errors.New("require a rabbitmq config")
	}

	if log == nil {
		return nil, errors.New("require a logging client")
	}

	client := Client{
		config:   conf,
		cbPool:   cbPool,
		ldapPool: ldapPool,
		mqConfig: mqConfig,
		log:      log,
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return &client, nil
}

// EOF
