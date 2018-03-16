package worker

import (
	"time"

	"github.com/zk101/nixie/app/loadtest/data"
)

// Worker is the actual work code
func (c *Client) Worker(local *data.Local) {
	c.WG.Add(1)
	config := c.testProfile.GetConfig()
	waitTimeBase := time.Second / time.Duration(config.RPS)
	currentDuration := time.Duration(0)
	count := uint32(0)
	testList := c.testProfile.GetTests()

	if err := c.testProfile.Start(local); err != nil {
		c.log.Sugar().Errorw("worker error", "error", err.Error(), "worker_id", local.WorkerID)
		c.runFlag = false
	}

	for c.runFlag {
		for _, test := range *testList {
			for x := uint(0); x < test.Count; x++ {
				start := time.Now()
				status := true

				if err := test.Callback(local); err != nil {
					c.log.Sugar().Errorw("worker error", "error", err.Error(), "worker_id", local.WorkerID)
					status = false
				}
				count++

				c.prometheus.ObserveReqDuration(&start, test.ID)
				c.prometheus.IncReqCount(local.WorkerID, test.ID, status)

				currentDuration += time.Since(start)
				controlDuration := time.Duration(count) * waitTimeBase

				if currentDuration <= controlDuration {
					time.Sleep(controlDuration - currentDuration)
					currentDuration = controlDuration
				}

				if count >= config.RPS {
					if currentDuration > controlDuration {
						c.log.Sugar().Errorw("worker error", "error", "not keeping up with requested rps", "worker_id", local.WorkerID, "rps", config.RPS, "current_duration", currentDuration.String())
					}

					currentDuration = time.Duration(0)
					count = uint32(0)
				}
			}
		}
		if err := c.testProfile.Reset(local); err != nil {
			c.log.Sugar().Errorw("worker error", "error", err.Error(), "worker_id", local.WorkerID)
			c.runFlag = false
		}
	}

	if err := c.testProfile.Stop(local); err != nil {
		c.log.Sugar().Errorw("worker error", "error", err.Error(), "worker_id", local.WorkerID)
	}

	c.WG.Done()
}

// EOF
