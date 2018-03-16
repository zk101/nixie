package chat

import (
	"errors"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/models/protobuf/chat/chatctl"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	"github.com/zk101/nixie/proto/chat"
	mChat "github.com/zk101/nixie/proto/chat"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Null sends a ChatMsg of type Null.  This is a chat login message
func Null(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	chatNull := chatctl.New()
	chatNull.Type = chat.ChatCtlType_CTL_TYPE_NULL
	chatNull.Userid = local.User

	dataNull, err := chatNull.Pack()
	msgTrack := uuid.NewV4().String()

	msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_CHAT_CTL, pbmsgwrap.MsgSec_SEC_SIGN, local.Key, local.Sign, local.Cipher, msgTrack, &dataNull)
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

	chatNullReply := chatctl.New()

	if err := chatNullReply.Unpack(wsmsg.GetMsgData()); err != nil {
		return err
	}

	if chatNullReply.Status != mChat.ChatCtlStatus_CTL_STATUS_OKAY {
		return errors.New("chat null returned failed status")
	}

	local.Peers[local.WorkerID].Online = true

	return nil
}

// EOF
