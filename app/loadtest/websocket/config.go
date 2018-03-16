package websocket

// Config holds websocket data
type Config struct {
	Host            string
	ReadBufferSize  int
	WriteBufferSize int
	LoopTime        int
}

// DefaultConfig returns a ready Config struct
func DefaultConfig() Config {
	return Config{
		Host:            "ws://localhost:10000",
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		LoopTime:        50,
	}
}

// EOF
