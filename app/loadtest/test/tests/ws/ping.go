package ws

import (
	"errors"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	mPingpong "github.com/zk101/nixie/models/protobuf/ws/pingpong"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Ping sends a ping and expects a response - Signed
func Ping(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	msgPing, sequence, err := mPingpong.PackClient()
	if err != nil {
		return err
	}

	msgTrack := uuid.NewV4().String()

	msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_PINGPONG, pbmsgwrap.MsgSec_SEC_SIGN, local.Key, local.Sign, local.Cipher, msgTrack, &msgPing)
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

	wsmsg, err := mMsgwrap.Unpack(&msgData.MsgData, local.Sign, local.Cipher)
	if err != nil {
		return err
	}

	msg, err := mPingpong.Unpack(wsmsg.GetMsgData())
	if err != nil {
		return err
	}

	if msg.GetSequence() != sequence {
		return errors.New("sequence does not match")
	}

	return nil
}

// EOF
