package gopool

import (
	"errors"
	"time"
)

// Schedule schedules task to be executed over pool's workers.
func (c *Client) Schedule(task func()) error {
	return c.schedule(task, time.After(time.Millisecond*time.Duration(c.config.ScheduleTimeout)))
}

// ScheduleTimeout schedules a task and overrides the timeout, used for net accept
func (c *Client) ScheduleTimeout(timeout time.Duration, task func()) error {
	return c.schedule(task, time.After(timeout))
}

// schedule does the actual scheduling
func (c *Client) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return errors.New("timeout scheduling task")
	case c.work <- task:
		return nil
	case c.register <- struct{}{}:
		go c.worker(task)
		return nil
	}
}

// EOF
