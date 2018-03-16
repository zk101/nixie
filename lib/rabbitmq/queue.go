package rabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

// QueueDeclare declares a Worker Queue of given name and returns the queue struct pointer of an error
func (c *Client) QueueDeclare(name string, durable bool) (*amqp.Queue, error) {
	if name == "" {
		return nil, errors.New("queue must have a name")
	}

	if c.Channel == nil {
		return nil, errors.New("channel not running")
	}

	queue, err := c.Channel.QueueDeclare(name, durable, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Channel.Qos(1, 0, false); err != nil {
		return nil, err
	}

	return &queue, nil
}

// QueueDelete removes a queue from rabbitmq
func (c *Client) QueueDelete(name string) (int, error) {
	if name == "" {
		return 0, errors.New("queue must have a name")
	}

	if c.Channel == nil {
		return 0, errors.New("channel not running")
	}

	count, err := c.Channel.QueueDelete(name, false, false, false)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// EOF
