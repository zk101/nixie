package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zk101/nixie/app/loadtest/prometheus"
	"github.com/zk101/nixie/lib/logging"
)

// Client holds operational vars forwebsocket connections
type Client struct {
	config     *Config
	log        *logging.Client
	prometheus *prometheus.Config
	conn       *websocket.Conn
	run        bool
	wg         *sync.WaitGroup
	writeChan  chan writeData
	quitChan   chan bool
	trackMap   map[string]Track
	autoMap    map[string]func([]byte) error
}

// NewClient returns a pointer to a new client
func NewClient(conf *Config, log *logging.Client, prometheus *prometheus.Config, autoMap map[string]func([]byte) error) *Client {
	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	client := Client{
		config:     conf,
		log:        log,
		prometheus: prometheus,
		run:        true,
		wg:         &sync.WaitGroup{},
		writeChan:  make(chan writeData, 10),
		quitChan:   make(chan bool, 2),
		trackMap:   make(map[string]Track, 0),
		autoMap:    autoMap,
	}

	return &client
}

// EOF
