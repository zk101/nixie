package prometheus

import (
	"time"

	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	ldappool "github.com/zk101/nixie/lib/ldap/pool"
	"github.com/zk101/nixie/lib/signal"
)

// Worker runs and polls metrics on the varies pools
func (c *Config) Worker(sig *signal.Config, cbpool *cbpool.Client, ldappool *ldappool.Client) {
	for sig.Run {
		c.SetLDAPQueueCount(ldappool.GetQueueCount())
		c.SetLDAPWorkerCount(ldappool.GetWorkerCount())
		c.SetCBQueueCount(cbpool.GetQueueCount())
		c.SetCBWorkerCount(cbpool.GetWorkerCount())

		time.Sleep(100 * time.Millisecond)
	}
}

// EOF
