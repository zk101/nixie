package user

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/zk101/nixie/app/ws/asynctx"
	"github.com/zk101/nixie/app/ws/prometheus"
	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/logging"
)

// Data stores per connection data and operational variables
type Data struct {
	io         sync.Mutex
	conn       io.ReadWriteCloser
	id         uint64
	key        string
	sign       string
	cipher     string
	timeStart  int64
	timeCheck  int64
	latency    int64
	serviceID  string
	log        *logging.Client
	prometheus *prometheus.Config
	asynctx    *asynctx.Client
	cbPool     *cbpool.Client
}

// NewData returns a configured Data struct pointer
func NewData(serviceID string, conn net.Conn, log *logging.Client, prometheus *prometheus.Config, asynctx *asynctx.Client, cbPool *cbpool.Client) *Data {
	return &Data{
		conn:       conn,
		timeStart:  time.Now().Unix(),
		serviceID:  serviceID,
		log:        log,
		prometheus: prometheus,
		asynctx:    asynctx,
		cbPool:     cbPool,
	}
}

// GetID returns the id field
func (d *Data) GetID() uint64 {
	return d.id
}

// SetID sets the id field
func (d *Data) SetID(id uint64) {
	d.id = id
}

// GetStartTime returns the start time
func (d *Data) GetStartTime() int64 {
	return d.timeStart
}

// GetKey returns the key field
func (d *Data) GetKey() string {
	return d.key
}

// EOF
