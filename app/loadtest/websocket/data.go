package websocket

// ReadData holds data from a websocket read
type ReadData struct {
	MsgType int
	MsgData []byte
	Err     error
}

// writeData holds data for writing to a socket
type writeData struct {
	msgType int
	msgData []byte
	errChan chan error
}

// EOF
