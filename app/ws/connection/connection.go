package connection

import (
	"net"

	"github.com/zk101/nixie/app/ws/connection/user"
)

// Register adds a new connection to the pool
func (c *Client) Register(conn net.Conn) *user.Data {
	user := user.NewData(c.serviceID, conn, c.log, c.prometheus, c.asynctx, c.cbPool)

	c.mu.Lock()
	user.SetID(c.userSeq)
	c.userMap[c.userSeq] = user

	c.userSeq++
	if c.userSeq >= 18446744073709551615 {
		c.userSeq = 0
	}
	c.mu.Unlock()

	c.prometheus.SetConnCount(len(c.userMap))
	c.prometheus.IncConnTotal()

	return user
}

// Remove discards a Connection
func (c *Client) Remove(user *user.Data) {
	c.mu.Lock()
	if _, has := c.userMap[user.GetID()]; has == true {
		delete(c.userMap, user.GetID())
	}
	c.mu.Unlock()

	c.prometheus.SetConnCount(len(c.userMap))
	c.prometheus.ObserveConnDuration(user.GetStartTime())
}

// EOF
