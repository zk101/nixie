package worker

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
	"github.com/zk101/nixie/models/couchbase/auth/presence"
	"github.com/zk101/nixie/models/ldap/chat/userchat"
	"github.com/zk101/nixie/models/protobuf/chat/chatmsg"
)

// msgChat deals with Post Messages type Test, Binary and Receipt
func (c *Client) msgChat(task *amqp.Delivery, msg *chatmsg.Model) error {

	// Check Friends Group

	fmt.Printf("%+v\n", msg)

	return nil
	//return c.SendChatMsg(task, msg)
}

// msgFriend deals with Post Messages type Friend
func (c *Client) msgFriend(task *amqp.Delivery, msg *chatmsg.Model) error {
	fromPresence := presence.New()
	if err := fromPresence.Fetch(c.cbPool, msg.From, true); err != nil {
		return err
	}

	toPresence := presence.New()
	if err := toPresence.Fetch(c.cbPool, msg.To, true); err != nil {
		return err
	}

	fromChatFriends := fromPresence.GetChatfriends()
	if value, ok := fromChatFriends[toPresence.GetDn()]; ok {
		if value == true {
			return nil
		}
	} else {
		fromChatData := userchat.New()
		fromChatData.DN = fromPresence.GetDn()
		if err := fromChatData.Fetch(c.ldapPool); err != nil {
			return err
		}

		fromChatData.Friends[toPresence.GetDn()] = false
		if err := fromChatData.Edit(c.ldapPool); err != nil {
			return err
		}

		// ChatFriends seems to be nil sometimes, which is a problem, not sure why its sometimes nil
		// New model should hopefully fix that
		if fromPresence.GetChatfriends() == nil {
			return errors.New("fromPresence is nil, not good")
		}

		fromChatFriends[toPresence.GetDn()] = false
		fromPresence.SetChatfriends(fromChatFriends)
		if err := fromPresence.Edit(c.cbPool); err != nil {
			return err
		}

		return c.SendChatMsg(task, msg)
	}

	return nil
}

// EOF
