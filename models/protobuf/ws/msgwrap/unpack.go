package msgwrap

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/lib/security"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Unpack returns a MsgWrap pointer
func Unpack(data *[]byte, msgSign, msgCipher string) (*pbmsgwrap.MsgWrap, error) {
	wsmsg := pbmsgwrap.MsgWrap{}
	if err := proto.Unmarshal(*data, &wsmsg); err != nil {
		return nil, err
	}

	switch wsmsg.GetMsgSec() {
	case pbmsgwrap.MsgSec_SEC_SIGN:
		if msgSign == "" {
			return nil, errors.New("requested msgsec requires signing secret")
		}
		signTest := security.CalcPBhash(msgSign, wsmsg.GetMsgKey(), wsmsg.GetMsgTrack(), int(wsmsg.GetMsgType()), int(wsmsg.GetMsgSec()), wsmsg.GetMsgDate(), wsmsg.GetMsgData())
		if wsmsg.GetMsgHash() != signTest {
			return nil, errors.New("hash compare failed")
		}

	case pbmsgwrap.MsgSec_SEC_CIPHER:
		if msgSign == "" {
			return nil, errors.New("requested msgsec requires signing secret")
		}
		if msgCipher == "" {
			return nil, errors.New("requested msgsec requires cipher secret")
		}
		signTest := security.CalcPBhash(msgSign, wsmsg.GetMsgKey(), wsmsg.GetMsgTrack(), int(wsmsg.GetMsgType()), int(wsmsg.GetMsgSec()), wsmsg.GetMsgDate(), wsmsg.GetMsgData())
		if wsmsg.GetMsgHash() != signTest {
			return nil, errors.New("hash compare failed")
		}

		pt, err := security.Decrypt(msgCipher, wsmsg.GetMsgData())
		if err != nil {
			return nil, err
		}
		wsmsg.MsgData = pt
	}

	return &wsmsg, nil
}

// EOF
