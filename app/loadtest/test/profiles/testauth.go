package profiles

import (
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/app/loadtest/test"
	"github.com/zk101/nixie/app/loadtest/test/tests/auth"
)

// TestAuth holds operational values and conforms to the profile interface
type TestAuth struct {
	config *test.Config
	tests  []test.Data
}

// Setup configures TestAuth
func (c *TestAuth) Setup(conf *test.Config) error {
	if conf == nil {
		c := test.DefaultConfig()
		conf = &c
	}
	c.config = conf

	c.tests = append(c.tests, test.Data{ID: "AuthHealth", Callback: auth.Health, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthRegister", Callback: auth.Register, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthLogin", Callback: auth.Login, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthRefresh", Callback: auth.Refresh, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthLogout", Callback: auth.Logout, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthLogin", Callback: auth.Login, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthRefresh", Callback: auth.Refresh, Count: 1})
	c.tests = append(c.tests, test.Data{ID: "AuthDeregister", Callback: auth.Deregister, Count: 1})

	return nil
}

// GetConfig returns the test Config ptr
func (c *TestAuth) GetConfig() *test.Config {
	return c.config
}

// GetTests returns the test profile
func (c *TestAuth) GetTests() *[]test.Data {
	return &c.tests
}

// Start sets up a local data instance with running services for use during testing
func (c *TestAuth) Start(local *data.Local) error {
	local.User = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Pass = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Name = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Peers[local.WorkerID].User = local.User

	local.AuthBase = c.config.TargetAuth

	return nil
}

// Stop shuts down running service in a local data instance
func (c *TestAuth) Stop(local *data.Local) error {
	return nil
}

// Reset is run at the end of end test cycle to allow for updates/resets to be made durubg a running test
func (c *TestAuth) Reset(local *data.Local) error {
	local.User = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Pass = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Name = strings.Replace(uuid.NewV4().String(), "-", "", -1)
	local.Peers[local.WorkerID].User = local.User

	return nil
}

// EOF
