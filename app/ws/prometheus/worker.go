package prometheus

import (
	"time"

	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	"github.com/zk101/nixie/lib/gopool"
	"github.com/zk101/nixie/lib/signal"
)

// Worker runs and polls metrics on the varies pools
func (c *Config) Worker(sig *signal.Config, cbPool *cbpool.Client, goPool *gopool.Client) {
	for sig.Run {
		c.SetCBQueueCount(cbPool.GetQueueCount())
		c.SetCBWorkerCount(cbPool.GetWorkerCount())
		c.SetGoPoolQueueCount(goPool.GetQueueCount())
		c.SetGoPoolWorkerCount(goPool.GetWorkerCount())

		time.Sleep(100 * time.Millisecond)
	}
}

// EOF
