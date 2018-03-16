package worker

import (
	"time"

	"github.com/streadway/amqp"
	"github.com/zk101/nixie/lib/rabbitmq"
	"github.com/zk101/nixie/models/sql/telemetry/telemetry"
)

// Worker is the actual work code
func (c *Client) worker(workerID string) {
	c.wg.Add(1)

	timeout := make(chan bool, 1)
	loopControl := true

	go func() {
		for {
			if c.run == false {
				timeout <- true
				break
			}
			time.Sleep(time.Second)
		}
	}()

	mq, err := rabbitmq.NewClient(c.mqConfig)
	if err != nil {
		c.log.Sugar().Errorw("rabbitmq connect failed", "worker_id", workerID, "error", err.Error())
		<-c.register
		c.wg.Done()
		time.Sleep(time.Millisecond * 5)
		return
	}
	defer mq.Close()

	msgs, err := mq.Consume("async_ff_telemetry_queue", "", false, false, false, false, nil)
	if err != nil {
		c.log.Sugar().Errorw("queue declare failed", "worker_id", workerID, "error", err.Error())
		<-c.register
		c.wg.Done()
		time.Sleep(time.Millisecond * 5)
		return
	}

	c.log.Sugar().Debugw("worker started", "worker_id", workerID)

	for loopControl {
		select {
		case <-timeout:
			loopControl = false
		case task := <-msgs:
			start := time.Now().Unix()
			state := "good"

			if len(task.Body) < 1 {
				state = "badmq"
				loopControl = false
				c.log.Sugar().Errorw("message queue failed", "worker_id", workerID)
			} else {
				if err := c.process(&task); err != nil {
					state = "bad"
					c.log.Sugar().Errorw("task process failed", "worker_id", workerID, "error", err.Error())
				}
			}

			task.Ack(false)

			c.log.Sugar().Debugw("telemetry processed", "worker_id", workerID)

			c.prometheus.IncTelemetryProcessCount(state)
			c.prometheus.ObserveProcessDuration(start)
		}
	}

	c.log.Sugar().Debugw("worker stopped", "worker_id", workerID)

	<-c.register
	c.wg.Done()
}

// process runs an acutal task
func (c *Client) process(task *amqp.Delivery) error {
	telemetryModel := telemetry.Model{}

	if err := telemetryModel.Unpack(task.Body); err != nil {
		return err
	}

	return telemetryModel.Create(c.sql)
}

// EOF
