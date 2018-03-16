package httpd

import (
	"errors"
	"net/http"

	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/mariadb"
	"github.com/zk101/nixie/lib/rabbitmq"
)

// Client adds operation struct to http which supports additional pointers
type Client struct {
	config   *httpd.Config
	mqConfig *rabbitmq.Config
	sqlMan   *mariadb.Config
	sqlRW    *mariadb.Config
	log      *logging.Client
	server   *http.Server
}

// NewClient returns a configured Client pointer
func NewClient(conf *httpd.Config, mqConfig *rabbitmq.Config, sqlMan *mariadb.Config, sqlRW *mariadb.Config, log *logging.Client) (*Client, error) {
	if conf == nil {
		return nil, errors.New("require an http config")
	}

	if mqConfig == nil {
		return nil, errors.New("require a rabbitmq config")
	}

	if sqlMan == nil {
		return nil, errors.New("require a sql manager config")
	}

	if sqlRW == nil {
		return nil, errors.New("require a sql rw config")
	}

	if log == nil {
		return nil, errors.New("require a logging client")
	}

	client := Client{
		config:   conf,
		mqConfig: mqConfig,
		sqlMan:   sqlMan,
		sqlRW:    sqlRW,
		log:      log,
	}

	if err := client.Start(); err != nil {
		return nil, err
	}

	return &client, nil
}

// EOF
