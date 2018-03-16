package config

import (
	"github.com/zk101/nixie/app/loadtest/test"
	"github.com/zk101/nixie/app/loadtest/websocket"
	"github.com/zk101/nixie/lib/consul"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/logging"
	"github.com/zk101/nixie/lib/signal"
)

// Config holds core configuration
type Config struct {
	Consul    consul.Config
	Controls  Controls
	HTTPD     httpd.Config
	Log       logging.Config
	Signal    signal.Config
	Test      test.Config
	WebSocket websocket.Config
}

// NewConfig creates a Config struct and returns a pointer
func NewConfig() *Config {
	return &Config{
		Consul:    consul.DefaultConfig(),
		Controls:  DefaultControls(),
		HTTPD:     httpd.DefaultConfig(),
		Log:       logging.DefaultConfig(),
		Signal:    signal.DefaultConfig(),
		Test:      test.DefaultConfig(),
		WebSocket: websocket.DefaultConfig(),
	}
}

// EOF
