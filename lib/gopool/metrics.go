package gopool

// GetWorkerCount returns the current number of workers
func (c *Client) GetWorkerCount() int {
	return len(c.register)
}

// GetQueueCount returns the current number of tasks on the queue
func (c *Client) GetQueueCount() int {
	return len(c.work)
}

// EOF
