package connection

import (
	"errors"
	"sync"

	"github.com/zk101/nixie/app/ws/asynctx"
	"github.com/zk101/nixie/app/ws/connection/user"
	"github.com/zk101/nixie/app/ws/prometheus"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/logging"
)

// Client holds operatonal values
type Client struct {
	mu         sync.RWMutex
	userSeq    uint64
	userMap    map[uint64]*user.Data
	serviceID  string
	log        *logging.Client
	prometheus *prometheus.Config
	asynctx    *asynctx.Client
	cbPool     *cbpool.Client
}

// NewClient sets up a new Client and returns a pointer
func NewClient(serviceID string, log *logging.Client, prometheus *prometheus.Config, asynctx *asynctx.Client, cbPool *cbpool.Client) (*Client, error) {
	if serviceID == "" {
		return nil, errors.New("require a valid service id")
	}

	if log == nil {
		return nil, errors.New("require a valid logging client")
	}

	if prometheus == nil {
		return nil, errors.New("require a valid prometheus config")
	}

	if asynctx == nil {
		return nil, errors.New("require a valid asynctx client")
	}

	if cbPool == nil {
		return nil, errors.New("require couchbase pool")
	}

	return &Client{
		userMap:    make(map[uint64]*user.Data),
		serviceID:  serviceID,
		log:        log,
		prometheus: prometheus,
		asynctx:    asynctx,
		cbPool:     cbPool,
	}, nil
}

// GetUserMap returns a pointer to the userMap
func (c *Client) GetUserMap() map[uint64]*user.Data {
	return c.userMap
}

// EOF
