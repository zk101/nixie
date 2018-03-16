package asyncrx

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/zk101/nixie/lib/rabbitmq"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// worker is the worker code
func (c *Client) worker() {
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
		c.log.Sugar().Errorw("rabbitmq connect failed", "process", "asyncrx", "error", err.Error())
		<-c.register
		c.wg.Done()
		return
	}
	defer mq.Close()

	msgs, err := mq.Consume("async_rpc_reply_"+c.serviceID, "", false, false, false, false, nil)
	if err != nil {
		c.log.Sugar().Errorw("queue declare failed", "process", "asyncrx", "error", err.Error())
		<-c.register
		c.wg.Done()
		time.Sleep(time.Millisecond * 5)
		return
	}

	for loopControl {
		select {
		case <-timeout:
			loopControl = false
		case task := <-msgs:
			task.Ack(false)
			queue := "async_rpc_reply_" + c.serviceID

			connParts := strings.Split(task.CorrelationId, "_")
			if len(connParts) != 4 {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asyncrx", "error", errors.New("correlation parts fo not equal four"))
				c.prometheus.IncAsyncrxDeliveryCount(queue, "bad")
				continue
			}

			connID, err := strconv.ParseUint(connParts[2], 10, 64)
			if err != nil {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asyncrx", "error", err)
				c.prometheus.IncAsyncrxDeliveryCount(queue, "bad")
				continue
			}

			msgType, err := strconv.Atoi(connParts[3])
			if err != nil {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asyncrx", "error", err)
				c.prometheus.IncAsyncrxDeliveryCount(queue, "bad")
				continue
			}

			if _, ok := c.userMap[connID]; ok == false {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asyncrx", "error", errors.New("userdata lookup failed"))
				c.prometheus.IncAsyncrxDeliveryCount(queue, "bad")
				continue
			}

			if c.userMap[connID].GetKey() != connParts[0] {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asyncrx", "error", errors.New("user key does not match connection key"))
				c.prometheus.IncAsyncrxDeliveryCount(queue, "bad")
				continue
			}

			if err := c.userMap[connID].Write(pbmsgwrap.MsgType(msgType), pbmsgwrap.MsgSec_SEC_SIGN, connParts[1], &task.Body); err != nil {
				loopControl = false
				c.log.Sugar().Errorw("task processing failed", "process", "asyncrx", "error", err)
				c.prometheus.IncAsyncrxDeliveryCount(queue, "bad")
				continue
			}

			lastWork = time.Now().Add(time.Second * time.Duration(c.config.WorkerExpiry))
			c.prometheus.IncAsyncrxDeliveryCount(queue, "good")
		}
	}

	<-c.register
	c.wg.Done()
}

// EOF
