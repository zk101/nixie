package asynctx

import (
	"time"

	"github.com/zk101/nixie/lib/rabbitmq"
)

// worker is the worker code
func (c *Client) worker(task task, initialStart bool) {
	c.wg.Add(1)

	lastWork := time.Now().Add(time.Second * time.Duration(c.config.WorkerExpiry))
	timeout := make(chan bool, 1)
	loopControl := true

	go func() {
		for {
			if time.Now().After(lastWork) {
				timeout <- true
				break
			}
			if c.run == false {
				timeout <- true
				break
			}
			time.Sleep(time.Second)
		}
	}()

	mq, err := rabbitmq.NewClient(c.mqConfig)
	if err != nil {
		loopControl = false
		initialStart = false
		c.log.Sugar().Errorw("rabbitmq connect failed", "process", "asynctx", "error", err.Error())
		<-c.register
		c.wg.Done()
		return
	}
	defer mq.Close()

	if initialStart == false {
		if err := mq.Publish("", task.queue, false, false, task.msg); err != nil {
			loopControl = false
			c.log.Sugar().Errorw("task processing failed", "process", "asynctx", "error", err.Error())
		}
	}

	for loopControl {
		select {
		case <-timeout:
			loopControl = false
		case task := <-c.work:
			if err := mq.Publish("", task.queue, false, false, task.msg); err != nil {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asynctx", "error", err.Error())
			}
			lastWork = time.Now().Add(time.Second * time.Duration(c.config.WorkerExpiry))
		}
	}

	<-c.register
	c.wg.Done()
}

// EOF
