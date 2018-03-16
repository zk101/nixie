package lib

import "github.com/zk101/nixie/lib/rabbitmq"

// SetupQueues connects to rabbitmq and creates all queues used by other processes
func (c *Core) SetupQueues() error {
	queues := [...]string{"async_rpc_chat_ctl_queue", "async_rpc_chat_msg_queue"}

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

	return nil
}

// EOF
