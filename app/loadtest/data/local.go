package data

import (
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/prometheus"
	"github.com/zk101/nixie/app/loadtest/websocket"
	"github.com/zk101/nixie/lib/logging"
)

// Local holds data required by the worker to run tests.  It also provides the ability to pass data into tests
type Local struct {
	WorkerID   string
	Key        string
	Sign       string
	Cipher     string
	User       string
	Pass       string
	Name       string
	AuthBase   string
	Peers      map[string]*Peer
	AutoAction map[string]func([]byte) error
	WSconf     *websocket.Config
	WS         *websocket.Client
	Log        *logging.Client
	Prometheus *prometheus.Config
}

// CreateLocal returns a ptr to a new Local struct
func CreateLocal() *Local {
	return &Local{
		Key:        uuid.NewV4().String(),
		Peers:      make(map[string]*Peer, 0),
		AutoAction: make(map[string]func([]byte) error, 0),
	}
}

// EOF
