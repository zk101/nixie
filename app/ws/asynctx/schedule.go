package asynctx

import (
	"errors"
	"strconv"
	"time"

	"github.com/zk101/nixie/lib/rabbitmq"
)

// Schedule schedules task to be executed over pool's workers.
func (c *Client) Schedule(connID uint64, key, track, queue string, persist bool, body []byte) error {
	if key == "" {
		return errors.New("key must not be empty")
	}

	if queue == "" {
		return errors.New("queue must not be empty")
	}

	msg, err := rabbitmq.CreateBodyRPC("application/octet-stream", key+"_"+track+"_"+strconv.FormatUint(connID, 10), "async_rpc_reply_"+c.serviceID, persist, body)
	if err != nil {
		return err
	}

	task := task{
		queue: queue,
		msg:   msg,
	}

	return c.schedule(task, time.After(time.Millisecond*time.Duration(c.config.ScheduleTimeout)))
}

// schedule does the actual scheduling
func (c *Client) schedule(task task, timeout <-chan time.Time) error {
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

// EOF
