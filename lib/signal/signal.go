package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Config holds signal data
type Config struct {
	Channel chan os.Signal
	Run     bool
}

// DefaultConfig initilises a SignalConf struct
func DefaultConfig() Config {
	c := Config{
		Channel: make(chan os.Signal, 1),
		Run:     true,
	}

	signal.Notify(c.Channel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	return c
}

// EOF
