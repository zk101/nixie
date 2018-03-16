package msgwrap

import (
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/lib/security"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Pack returns a []byte ready to send
func Pack(msgType pbmsgwrap.MsgType, msgSec pbmsgwrap.MsgSec, msgKey, msgSign, msgCipher, msgTrack string, msgData *[]byte) (*[]byte, error) {
	if msgKey == "" {
		return nil, errors.New("require a msgKey value")
	}

	if msgTrack == "" {
		return nil, errors.New("require a msgTrack value")
	}

	var msg []byte
	var err error
	unixTime := time.Now().Unix()

	msgwrap := &pbmsgwrap.MsgWrap{
		MsgType:  msgType,
		MsgSec:   msgSec,
		MsgTrack: msgTrack,
		MsgDate:  unixTime,
		MsgKey:   msgKey,
	}

	switch msgSec {
	case pbmsgwrap.MsgSec_SEC_SIGN:
		if msgSign == "" {
			return nil, errors.New("requested msgsec requires signing secret")
		}
		msgwrap.MsgHash = security.CalcPBhash(msgSign, msgKey, msgTrack, int(msgType), int(msgSec), unixTime, *msgData)
		msgwrap.MsgData = *msgData

	case pbmsgwrap.MsgSec_SEC_CIPHER:
		if msgSign == "" {
			return nil, errors.New("requested msgsec requires signing secret")
		}
		if msgCipher == "" {
			return nil, errors.New("requested msgsec requires cipher secret")
		}
		ct, err := security.Encrypt(msgCipher, *msgData)
		if err != nil {
			return nil, err
		}
		msgwrap.MsgHash = security.CalcPBhash(msgSign, msgKey, msgTrack, int(msgType), int(msgSec), unixTime, ct)
		msgwrap.MsgData = ct

	default:
		if msgType == pbmsgwrap.MsgType_MSG_NULL {
			if msgSign == "" {
				return nil, errors.New("null msg requires signing secret")
			}
			msgwrap.MsgHash = security.CalcPBhash(msgSign, msgKey, msgTrack, int(msgType), int(msgSec), unixTime, *msgData)
		}
		msgwrap.MsgData = *msgData
	}

	if msg, err = proto.Marshal(msgwrap); err != nil {
		return nil, err
	}

	return &msg, nil
}

// EOF
