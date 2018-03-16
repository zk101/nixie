package chat

import (
	"errors"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/models/protobuf/chat/chatmsg"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	mChat "github.com/zk101/nixie/proto/chat"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Friend sends a ChatMsg of type Friend.
func Friend(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	for key, value := range local.Peers {
		if local.WorkerID == key {
			continue
		}

		chatFriend := chatmsg.New()
		chatFriend.Type = mChat.ChatMsgType_MSG_TYPE_FRIEND
		chatFriend.To = value.User
		chatFriend.From = local.User

		dataFriend, err := chatFriend.Pack()
		msgTrack := uuid.NewV4().String()

		msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_CHAT_MSG, pbmsgwrap.MsgSec_SEC_SIGN, local.Key, local.Sign, local.Cipher, msgTrack, &dataFriend)
		if err != nil {
			return err
		}

		tracker := local.WS.AddTrack(msgTrack)

		if err := local.WS.Write(gwebsock.BinaryMessage, *msgWrap); err != nil {
			return err
		}

		msgData, err := tracker.ReadBlock()
		if err != nil {
			return err
		}
		//msgData := <-tracker.Reply

		wsmsg, err := mMsgwrap.Unpack(&msgData.MsgData, local.Sign, local.Cipher)
		if err != nil {
			return err
		}

		chatFriendReply := chatmsg.New()

		if err := chatFriendReply.Unpack(wsmsg.GetMsgData()); err != nil {
			return err
		}
	}

	local.Peers[local.WorkerID].Friends = true

	return nil
}

// EOF
