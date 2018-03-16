package profiles

import (
	"github.com/zk101/nixie/app/loadtest/lib"
	"github.com/zk101/nixie/app/loadtest/test"
)

// GetTestProfiles returns a map of profiles
func GetTestProfiles(core *lib.Core) map[string]test.Profile {
	testList := make(map[string]test.Profile, 0)

	testList["auth"] = &TestAuth{}

	testList["all"] = &TestAll{
		WSconfig:   &core.Config.WebSocket,
		Log:        core.Clients.Log,
		Prometheus: core.Clients.Prometheus,
	}

	testList["peer"] = &TestPeer{
		WSconfig:   &core.Config.WebSocket,
		Log:        core.Clients.Log,
		Prometheus: core.Clients.Prometheus,
	}

	return testList
}

// EOF
