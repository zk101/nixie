package lib

import "github.com/zk101/nixie/lib/rabbitmq"

// SetupQueues connects to rabbitmq and creates all queues used by other processes
func (c *Core) SetupQueues() error {
	mq, err := rabbitmq.NewClient(&c.Config.RabbitMQ)
	if err != nil {
		return err
	}
	defer mq.Close()
	if _, err := mq.QueueDeclare("async_ff_telemetry_queue", true); err != nil {
		return err
	}

	return nil
}

// EOF
