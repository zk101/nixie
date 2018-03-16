package ws

import (
	"errors"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	mSrvtime "github.com/zk101/nixie/models/protobuf/ws/srvtime"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Srvtime sends a Server Time request
func Srvtime(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	var msgSrvTime []byte
	msgTrack := uuid.NewV4().String()

	msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_SRVTIME, pbmsgwrap.MsgSec_SEC_SIGN, local.Key, local.Sign, local.Cipher, msgTrack, &msgSrvTime)
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

	msg, err := mMsgwrap.Unpack(&msgData.MsgData, local.Sign, local.Cipher)
	if err != nil {
		return err
	}

	_, err = mSrvtime.Unpack(msg.GetMsgData())
	if err != nil {
		return err
	}

	return nil
}

// EOF
