package profiles

import (
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/app/loadtest/prometheus"
	"github.com/zk101/nixie/app/loadtest/test"
	"github.com/zk101/nixie/app/loadtest/test/tests/auth"
	"github.com/zk101/nixie/app/loadtest/test/tests/chat"
	"github.com/zk101/nixie/app/loadtest/test/tests/telemetry"
	"github.com/zk101/nixie/app/loadtest/test/tests/ws"
	"github.com/zk101/nixie/app/loadtest/websocket"
	"github.com/zk101/nixie/lib/logging"
)

// TestAll holds operational values and conforms to the profile interface
type TestAll struct {
	config     *test.Config
	tests      []test.Data
	WSconfig   *websocket.Config
	Log        *logging.Client
	Prometheus *prometheus.Config
}

// Setup configures TestAll
func (c *TestAll) Setup(conf *test.Config) error {
	if conf == nil {
		c := test.DefaultConfig()
		conf = &c
	}
	c.config = conf

	c.tests = append(c.tests, test.Data{ID: "AuthHealth", Callback: auth.Health, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthRegister", Callback: auth.Register, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthLogin", Callback: auth.Login, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthRefresh", Callback: auth.Refresh, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "WSconnect", Callback: ws.Connect, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSnull", Callback: ws.Null, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSlatency", Callback: ws.Latency, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSsrvtime", Callback: ws.Srvtime, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSping", Callback: ws.Ping, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "WStelemetry", Callback: telemetry.Send, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "WSchatNull", Callback: chat.Null, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "WSclose", Callback: ws.Close, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "AuthLogout", Callback: auth.Logout, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthLogin", Callback: auth.Login, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthRefresh", Callback: auth.Refresh, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "WSconnect", Callback: ws.Connect, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSnull", Callback: ws.Null, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSchatNullcreated", Callback: chat.Null, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "WSclose", Callback: ws.Close, Count: 1})

	c.tests = append(c.tests, test.Data{ID: "AuthDeregister", Callback: auth.Deregister, Count: 1})

	return nil
}

// GetConfig returns the test Config ptr
func (c *TestAll) GetConfig() *test.Config {
	return c.config
}

// GetTests returns the test profile
func (c *TestAll) GetTests() *[]test.Data {
	return &c.tests
}

// Start sets up a local data instance with running services for use during testing
func (c *TestAll) Start(local *data.Local) error {
	local.User = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Pass = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Name = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Peers[local.WorkerID].User = local.User

	local.AuthBase = c.config.TargetAuth
	local.WSconf = c.WSconfig
	local.Log = c.Log
	local.Prometheus = c.Prometheus

	return nil
}

// Stop shuts down running service in a local data instance
func (c *TestAll) Stop(local *data.Local) error {
	return nil
}

// Reset is run at the end of end test cycle to allow for updates/resets to be made durubg a running test
func (c *TestAll) Reset(local *data.Local) error {
	local.User = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Pass = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Name = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Peers[local.WorkerID].User = local.User

	return nil
}

// EOF
