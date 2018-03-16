package pool

import "github.com/mavricknz/ldap"

// doTask
func (c *Client) doTask(task taskRequest) *taskReply {
	reply := taskReply{}

	if reply.err = c.scheduleWrapper(task); reply.err != nil {
		return &reply
	}

	reply = <-task.reply
	return &reply
}

// Add is a helper wrapper
func (c *Client) Add(req *ldap.AddRequest) error {
	tr := c.doTask(addTask(req))
	if tr.err != nil {
		return tr.err
	}
	return nil
}

// Bind is a helper wrapper
func (c *Client) Bind(username, password string) error {
	tr := c.doTask(bindTask(username, password))
	if tr.err != nil {
		return tr.err
	}
	return nil
}

// Delete is a helper wrapper
func (c *Client) Delete(delReq *ldap.DeleteRequest) error {
	tr := c.doTask(deleteTask(delReq))
	if tr.err != nil {
		return tr.err
	}
	return nil
}

// Modify is a helper wrapper
func (c *Client) Modify(req *ldap.ModifyRequest) error {
	tr := c.doTask(modTask(req))
	if tr.err != nil {
		return tr.err
	}
	return nil
}

// ModDn is a helper wrapper
func (c *Client) ModDn(req *ldap.ModDnRequest) error {
	tr := c.doTask(modDnTask(req))
	if tr.err != nil {
		return tr.err
	}
	return nil
}

// Search is a helper wrapper
func (c *Client) Search(req *ldap.SearchRequest) (*ldap.SearchResult, error) {
	tr := c.doTask(searchTask(req))
	if tr.err != nil {
		return nil, tr.err
	}
	return tr.sr, nil
}

// EOF
