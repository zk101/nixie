package worker

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
	"github.com/zk101/nixie/lib/rabbitmq"
	"github.com/zk101/nixie/models/couchbase/auth/presence"
	"github.com/zk101/nixie/models/protobuf/chat/chatctl"
	"github.com/zk101/nixie/models/protobuf/chat/chatmsg"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// SendChatCtl packs and transmits a ctl message
func (c *Client) SendChatCtl(task *amqp.Delivery, msg *chatctl.Model) error {
	msgData, err := msg.Pack()
	if err != nil {
		return err
	}

	msgBody, err := rabbitmq.CreateBodyRPCReply("application/octet-stream", task.CorrelationId+"_"+strconv.Itoa(int(pbmsgwrap.MsgType_value["MSG_CHAT_CTL"])), true, msgData)
	if err != nil {
		return err
	}

	return c.mqClient.Publish("", task.ReplyTo, false, false, msgBody)
}

// SendChatMsg packs and transmits a ctl message
func (c *Client) SendChatMsg(task *amqp.Delivery, msg *chatmsg.Model) error {
	msgData, err := msg.Pack()
	if err != nil {
		return err
	}

	connParts := strings.Split(task.CorrelationId, "_")
	if len(connParts) != 3 {
		return errors.New("correlation parts fo not equal three")
	}

	modelPresence := presence.New()
	if err := modelPresence.Fetch(c.cbPool, msg.To, true); err != nil {
		return errors.New("get model presence failed")
	}

	correlationID := fmt.Sprintf("%s_%s_%s_%s", modelPresence.GetKey(), connParts[1], strconv.FormatUint(modelPresence.GetConnectionid(), 10), strconv.Itoa(int(pbmsgwrap.MsgType_value["MSG_CHAT_MSG"])))

	msgBody, err := rabbitmq.CreateBodyRPCReply("application/octet-stream", correlationID, true, msgData)
	if err != nil {
		return err
	}

	return c.mqClient.Publish("", "async_rpc_reply_"+modelPresence.GetServerid(), false, false, msgBody)
}

// EOF
