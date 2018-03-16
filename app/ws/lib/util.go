package lib

// NetpollError is used to allow logging from within Netpoll
func (c *Core) NetpollError(err error) {
	c.Clients.Log.Sugar().Errorf("netpoll error: %s", err.Error())
}

// EOF
