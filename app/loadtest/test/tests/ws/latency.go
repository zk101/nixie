package ws

import (
	"errors"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	mLatency "github.com/zk101/nixie/models/protobuf/ws/latency"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Latency sends a ping and expects a response
func Latency(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	msgLatency, err := mLatency.PackClient()
	if err != nil {
		return err
	}

	msgTrack := uuid.NewV4().String()

	msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_LATENCY, pbmsgwrap.MsgSec_SEC_SIGN, local.Key, local.Sign, local.Cipher, msgTrack, &msgLatency)
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

	_, err = mMsgwrap.Unpack(&msgData.MsgData, local.Sign, local.Cipher)
	if err != nil {
		return err
	}

	return nil
}

// EOF
