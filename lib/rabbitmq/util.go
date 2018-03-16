package rabbitmq

// Close hides amqp channel close and implements its own (which calls the channels close)
func (c *Client) Close() error {
	if c.Channel != nil {
		if err := c.Channel.Close(); err != nil {
			if c.conn != nil {
				if err := c.conn.Close(); err != nil {
					return err
				}
			}
			return err
		}
	}
	c.Channel = nil

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}
	c.conn = nil

	return nil
}

// EOF
