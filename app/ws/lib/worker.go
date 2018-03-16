package lib

import (
	"time"
)

// worker is a go routine that runs core application things...bit vague
func (c *Core) worker() {
	c.wg.Add(1)

	checkServiceCounter := 0

	for c.run {
		if checkServiceCounter >= 15 {
			if found := c.checkServiceExists(); found == false {
				if err := c.doConsulReg(); err != nil {
					c.Clients.Log.Sugar().Errorw("service re-registration failed", "service_id", c.Clients.ServiceID, "error", err.Error())
				} else {
					c.Clients.Log.Sugar().Debugw("service re-registered", "service_id", c.Clients.ServiceID)
				}
			}

			checkServiceCounter = 0
		}

		time.Sleep(time.Second)
		checkServiceCounter++
	}

	c.wg.Done()
}

// checkServiceExists looks at consul for its service id and returns true if the service appears there.
func (c *Core) checkServiceExists() bool {
	services, _, err := c.Clients.Consul.HealthService(c.Config.Controls.ServiceName, "", false, nil)
	if err != nil {
		c.Clients.Log.Sugar().Errorw("service lookup failed", "service_id", c.Clients.ServiceID, "error", err.Error())
		return false
	}

	for x := range services {
		for y := range services[x].Checks {
			if services[x].Checks[y].ServiceID == c.Clients.ServiceID {
				return true
			}
		}
	}

	return false
}

// EOF
