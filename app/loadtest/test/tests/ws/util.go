package ws

import (
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/app/loadtest/websocket"
)

// Connect connects a websocket
func Connect(local *data.Local) error {
	local.WS = websocket.NewClient(local.WSconf, local.Log, local.Prometheus, local.AutoAction)

	return local.WS.Start()
}

// Close shuts a websocket down
func Close(local *data.Local) error {
	return local.WS.Stop()
}

// EOF
