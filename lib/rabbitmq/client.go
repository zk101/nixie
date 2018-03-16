package rabbitmq

import (
	"github.com/streadway/amqp"
)

// Client holds working data
type Client struct {
	*amqp.Channel
	config *Config
	conn   *amqp.Connection
}

// NewClient returns a Client pointer
func NewClient(conf *Config) (*Client, error) {
	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	conn, err := amqp.Dial(conf.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{ch, conf, conn}, nil
}

// EOF
