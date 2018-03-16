package chat

import (
	"errors"
	"time"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/models/protobuf/chat/chatctl"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	"github.com/zk101/nixie/proto/chat"
	mChat "github.com/zk101/nixie/proto/chat"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Search sends a search request
func Search(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	for key, value := range local.Peers {
		if local.WorkerID == key {
			continue
		}

		loopBreak := 0
		for {
			if value.Online == false {
				if loopBreak > 2000 {
					return errors.New("chat search gave up waiting for peers")
				}

				time.Sleep(time.Millisecond)
				loopBreak++
				continue
			}

			break
		}

		chatSearch := chatctl.New()
		chatSearch.Type = chat.ChatCtlType_CTL_TYPE_SEARCH
		chatSearch.Userid = value.User

		dataNull, err := chatSearch.Pack()
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

		chatSearchReply := chatctl.New()

		if err := chatSearchReply.Unpack(wsmsg.GetMsgData()); err != nil {
			return err
		}

		if chatSearchReply.Status != mChat.ChatCtlStatus_CTL_STATUS_OKAY {
			return errors.New("chat search returned failed status")
		}
	}

	return nil
}

// EOF
