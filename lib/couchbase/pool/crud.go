package pool

import (
	"github.com/couchbase/gocb"
)

// doTask
func (c *Client) doTask(task taskRequest) *taskReply {
	reply := taskReply{}

	if reply.err = c.scheduleWrapper(task); reply.err != nil {
		return &reply
	}

	reply = <-task.reply
	return &reply
}

// Get is a helper wrapper
func (c *Client) Get(key string, valuePtr interface{}) (gocb.Cas, error) {
	tr := c.doTask(getTask(key, valuePtr))
	if tr.err != nil {
		return 0, tr.err
	}

	return tr.cas, nil
}

// Insert is a helper wrapper
func (c *Client) Insert(key string, value interface{}, expiry uint32) (gocb.Cas, error) {
	tr := c.doTask(insertTask(key, value, expiry))
	if tr.err != nil {
		return 0, tr.err
	}

	return tr.cas, nil
}

// Remove is a helper wrapper
func (c *Client) Remove(key string, cas gocb.Cas) (gocb.Cas, error) {
	tr := c.doTask(removeTask(key, cas))
	if tr.err != nil {
		return 0, tr.err
	}

	return tr.cas, nil
}

// Replace is a helper wrapper
func (c *Client) Replace(key string, value interface{}, cas gocb.Cas, expiry uint32) (gocb.Cas, error) {
	tr := c.doTask(replaceTask(key, value, cas, expiry))
	if tr.err != nil {
		return 0, tr.err
	}

	return tr.cas, nil
}

// Touch is a helper wrapper
func (c *Client) Touch(key string, cas gocb.Cas, expiry uint32) (gocb.Cas, error) {
	tr := c.doTask(touchTask(key, cas, expiry))
	if tr.err != nil {
		return 0, tr.err
	}

	return tr.cas, nil
}

// Upsert is a helper wrapper
func (c *Client) Upsert(key string, value interface{}, expiry uint32) (gocb.Cas, error) {
	tr := c.doTask(upsertTask(key, value, expiry))
	if tr.err != nil {
		return 0, tr.err
	}

	return tr.cas, nil
}

// EOF
