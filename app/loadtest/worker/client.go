package worker

import (
	"sync"

	"github.com/zk101/nixie/app/loadtest/lib"
	"github.com/zk101/nixie/app/loadtest/prometheus"
	"github.com/zk101/nixie/app/loadtest/test"
	"github.com/zk101/nixie/lib/logging"
)

// Client holds worker config
type Client struct {
	testProfile test.Profile
	runFlag     bool
	WG          sync.WaitGroup
	log         *logging.Client
	prometheus  *prometheus.Config
}

// NewClient sets the test up
func NewClient(core *lib.Core, testProfile test.Profile) (*Client, error) {
	if err := testProfile.Setup(&core.Config.Test); err != nil {
		return nil, err
	}

	return &Client{
		testProfile: testProfile,
		runFlag:     true,
		WG:          sync.WaitGroup{},
		log:         core.Clients.Log,
		prometheus:  core.Clients.Prometheus,
	}, nil
}

// SwitchRunFlag toggles the controlFlag
func (c *Client) SwitchRunFlag() {
	c.runFlag = !c.runFlag
}

// EOF
