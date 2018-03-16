package lib

import (
	"errors"

	"github.com/zk101/nixie/lib/rabbitmq"
)

// SetupQueues connects to rabbitmq and creates all queues used by other processes
func (c *Core) SetupQueues() error {
	if c.Clients == nil {
		return errors.New("service id is located in clients which needs to be setup first")
	}

	if c.Clients.ServiceID == "" {
		return errors.New("require service id to exist")
	}

	queues := [...]string{"async_ff_telemetry_queue", "async_rpc_chat_ctl_queue", "async_rpc_chat_post_queue"}

	mq, err := rabbitmq.NewClient(&c.Config.RabbitMQ)
	if err != nil {
		return err
	}
	defer mq.Close()

	for _, value := range queues {
		if _, err := mq.QueueDeclare(value, true); err != nil {
			return err
		}
	}

	if _, err := mq.QueueDeclare("async_rpc_reply_"+c.Clients.ServiceID, false); err != nil {
		return err
	}

	return nil
}

// DeleteQueues removes application specific queues (using autodelete creates larger issues due to worker pools)
func (c *Core) DeleteQueues() error {
	mq, err := rabbitmq.NewClient(&c.Config.RabbitMQ)
	if err != nil {
		return err
	}
	defer mq.Close()

	if _, err := mq.QueueDelete("async_rpc_reply_" + c.Clients.ServiceID); err != nil {
		return err
	}

	return nil
}

// EOF
