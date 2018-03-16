package worker

import (
	"errors"
	"time"

	"github.com/zk101/nixie/proto/chat"

	"github.com/streadway/amqp"
	"github.com/zk101/nixie/lib/rabbitmq"
	"github.com/zk101/nixie/models/protobuf/chat/chatctl"
	"github.com/zk101/nixie/models/protobuf/chat/chatmsg"
	pbChat "github.com/zk101/nixie/proto/chat"
)

// Worker is the actual work code
func (c *Client) worker(workerID string) {
	c.wg.Add(1)

	timeout := make(chan bool, 1)
	loopControl := true
	var err error

	go func() {
		for {
			if c.run == false {
				timeout <- true
				break
			}
			time.Sleep(time.Second)
		}
	}()

	c.mqClient, err = rabbitmq.NewClient(c.mqConfig)
	if err != nil {
		c.log.Sugar().Errorw("rabbitmq connect failed", "worker_id", workerID, "error", err.Error())
		<-c.register
		c.wg.Done()
		time.Sleep(time.Millisecond * 5)
		return
	}
	defer c.mqClient.Close()

	queueCtl, err := c.mqClient.Consume("async_rpc_chat_ctl_queue", "", false, false, false, false, nil)
	if err != nil {
		c.log.Sugar().Errorw("queue consume failed", "worker_id", workerID, "error", err.Error())
		<-c.register
		c.wg.Done()
		time.Sleep(time.Millisecond * 5)
		return
	}

	queueMsg, err := c.mqClient.Consume("async_rpc_chat_msg_queue", "", false, false, false, false, nil)
	if err != nil {
		c.log.Sugar().Errorw("queue consume failed", "worker_id", workerID, "error", err.Error())
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

		case task := <-queueMsg:
			if len(task.Body) < 1 {
				loopControl = false
				c.log.Sugar().Errorw("message queue failed", "worker_id", workerID)
				task.Ack(false)
				break
			}

			start := time.Now().Unix()
			state := "good"
			var category string
			var err error

			if category, err = c.processMsg(&task); err != nil {
				state = "bad"
				c.log.Sugar().Errorw("task process failed", "worker_id", workerID, "error", err.Error())
			}

			task.Ack(false)

			c.log.Sugar().Debugw("chat post processed", "worker_id", workerID)

			c.prometheus.IncChatProcessCount(state, category)
			c.prometheus.ObserveProcessDuration(start, state, category)

		case task := <-queueCtl:
			if len(task.Body) < 1 {
				loopControl = false
				c.log.Sugar().Errorw("message queue failed", "worker_id", workerID)
				task.Ack(false)
				break
			}

			start := time.Now().Unix()
			state := "good"
			var category string
			var err error

			if category, err = c.processCtl(&task); err != nil {
				state = "bad"
				c.log.Sugar().Errorw("task process failed", "worker_id", workerID, "error", err.Error())
			}

			task.Ack(false)

			c.log.Sugar().Debugw("chat control processed", "worker_id", workerID)

			c.prometheus.IncChatProcessCount(state, category)
			c.prometheus.ObserveProcessDuration(start, state, category)
		}
	}

	c.log.Sugar().Debugw("worker stopped", "worker_id", workerID)

	<-c.register
	c.wg.Done()
}

// processMsg processes ChatPost messages
func (c *Client) processMsg(task *amqp.Delivery) (string, error) {
	var err error
	msgData := chatmsg.New()

	err = msgData.Unpack(task.Body)
	if err != nil {
		return "Bad ChatMsg Unpack", err
	}

	switch msgData.Type {
	case pbChat.ChatMsgType_MSG_TYPE_TEXT:
		err = c.msgChat(task, msgData)

	case pbChat.ChatMsgType_MSG_TYPE_BINARY:
		err = c.msgChat(task, msgData)

	case pbChat.ChatMsgType_MSG_TYPE_RECEIPT:
		err = c.msgChat(task, msgData)

	case pbChat.ChatMsgType_MSG_TYPE_FRIEND:
		err = c.msgFriend(task, msgData)

	default:
		return "unsupported chatmsg message", errors.New("unsupported chatmsg message")
	}

	return pbChat.ChatMsgType_name[int32(msgData.Type)], err
}

// processCtl process ChatControl messages
func (c *Client) processCtl(task *amqp.Delivery) (string, error) {
	var err error
	ctlData := chatctl.New()

	err = ctlData.Unpack(task.Body)
	if err != nil {
		return "Bad ChatCtl Unpack", err
	}

	ctlData.Date = time.Now().Unix()
	ctlData.Status = chat.ChatCtlStatus_CTL_STATUS_OKAY

	switch ctlData.Type {
	case pbChat.ChatCtlType_CTL_TYPE_NULL:
		err = c.ctlNull(ctlData)

	case pbChat.ChatCtlType_CTL_TYPE_SEARCH:
		err = c.ctlSearch(ctlData)

	default:
		return "unsupported chatctl message", errors.New("unsupported chatctl message")
	}

	if err != nil {
		ctlData.Status = chat.ChatCtlStatus_CTL_STATUS_FAIL
	}

	return pbChat.ChatCtlType_name[int32(ctlData.Type)], c.SendChatCtl(task, ctlData)
}

// EOF
