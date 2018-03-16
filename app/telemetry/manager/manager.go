package manager

import (
	"fmt"
	"time"
)

// manager is the actual work code
func (c *Client) manager() {
	c.wg.Add(1)

	timeoutRunFlag := make(chan bool, 1)
	timeoutManager := make(chan bool, 1)
	countManager := 0
	loopControl := true

	go func() {
		for {
			if c.run == false {
				timeoutRunFlag <- true
				break
			}
			if countManager > 30 {
				timeoutManager <- true
				countManager = 0
			}
			countManager++
			time.Sleep(time.Second)
		}
	}()

	if err := c.process(); err != nil {
		c.log.Sugar().Errorw("run failed", "worker_id", "manager", "error", err.Error())
	}

	for loopControl {
		select {
		case <-timeoutRunFlag:
			loopControl = false
		case <-timeoutManager:
			if err := c.process(); err != nil {
				c.log.Sugar().Errorw("run failed", "worker_id", "manager", "error", err.Error())
			}
		}
	}

	c.wg.Done()
}

func telemetryTable(name string) string {
	table := "CREATE TABLE IF NOT EXISTS `%s` LIKE telemetry"

	return fmt.Sprintf(table, name)
}

func (c *Client) process() error {
	datetime := time.Now()

	for x := 0; x < 12; x++ {
		tableName := datetime.Add(time.Duration(x)*time.Hour).UTC().Format("200601021500") + "_telemetry"

		if _, err := c.sql.Exec(telemetryTable(tableName)); err != nil {
			c.log.Sugar().Errorw("manager table create failed", "worker_id", "manager", "error", err.Error())
		}
	}

	return nil
}

// EOF
