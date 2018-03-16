package websocket

import (
	"github.com/gorilla/websocket"
)

// connect connects to a websocket
func (c *Client) connect() error {
	if c.conn != nil {
		c.conn = nil
	}

	dialerConf := websocket.Dialer{
		ReadBufferSize:    c.config.ReadBufferSize,
		WriteBufferSize:   c.config.WriteBufferSize,
		EnableCompression: false,
	}

	var err error

	c.conn, _, err = dialerConf.Dial(c.config.Host, nil)
	if err != nil {
		return err
	}

	return nil
}

// close shuts down a running connection
func (c *Client) close() error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		c.conn = nil
		return err
	}

	c.conn = nil

	return nil
}

// Start sets up the websocket connections and starts the read and write workers
func (c *Client) Start() error {
	if err := c.connect(); err != nil {
		return err
	}

	go c.workerRead()
	go c.workerWrite()

	return nil
}

// Stop shuts down the two workers and closes the websocket connection
func (c *Client) Stop() error {
	errChan := make(chan error, 1)

	c.writeChan <- writeData{
		msgType: websocket.CloseMessage,
		msgData: websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		errChan: errChan,
	}

	if err := <-errChan; err != nil {
		return err
	}

	c.run = false
	c.quitChan <- true
	c.quitChan <- true
	c.wg.Wait()

	return c.close()
}

// Write sends a message
func (c *Client) Write(msgType int, msgData []byte) error {
	data := writeData{
		msgType: msgType,
		msgData: msgData,
		errChan: make(chan error, 1),
	}

	c.writeChan <- data

	return <-data.errChan
}

// EOF
