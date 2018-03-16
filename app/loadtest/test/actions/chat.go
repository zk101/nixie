package actions

import (
	"errors"
	"fmt"

	gwebsock "github.com/gorilla/websocket"
	"github.com/zk101/nixie/models/protobuf/chat/chatmsg"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	mChat "github.com/zk101/nixie/proto/chat"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// ChatMsg responds to ChatMsg messages
func (c *Config) ChatMsg(msg []byte) error {
	wsmsg, err := mMsgwrap.Unpack(&msg, c.Local.Sign, c.Local.Cipher)
	if err != nil {
		return err
	}

	chatMsg := chatmsg.New()
	if err := chatMsg.Unpack(wsmsg.GetMsgData()); err != nil {
		return err
	}

	switch chatMsg.Type {
	case mChat.ChatMsgType_MSG_TYPE_FRIEND:
		fmt.Println()
		fmt.Println("Auto Message Friend")
		fmt.Println()

		chatMsgReply := chatmsg.New()
		chatMsgReply.Type = mChat.ChatMsgType_MSG_TYPE_FRIEND
		chatMsgReply.To = chatMsg.From
		chatMsgReply.From = chatMsg.To

		dataFriend, err := chatMsgReply.Pack()

		msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_CHAT_MSG, pbmsgwrap.MsgSec_SEC_SIGN, c.Local.Key, c.Local.Sign, c.Local.Cipher, wsmsg.GetMsgTrack(), &dataFriend)
		if err != nil {
			return err
		}

		if err := c.Local.WS.Write(gwebsock.BinaryMessage, *msgWrap); err != nil {
			return err
		}

	default:
		return errors.New("autoaction chatmsg type unsupported")
	}

	return nil
}

// EOF
