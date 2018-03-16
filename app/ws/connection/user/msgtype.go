package user

import (
	"errors"
	"time"

	"github.com/zk101/nixie/lib/security"
	"github.com/zk101/nixie/models/couchbase/auth/presence"
	mLatency "github.com/zk101/nixie/models/protobuf/ws/latency"
	mPingpong "github.com/zk101/nixie/models/protobuf/ws/pingpong"
	mSrvtime "github.com/zk101/nixie/models/protobuf/ws/srvtime"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Ping Msg
func (d *Data) pingpingMsg(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_SIGN {
		return errors.New("required security level not set")
	}

	msg, err := mPingpong.PackServer(wsmsg.GetMsgData())
	if err != nil {
		return err
	}

	return d.Write(pbmsgwrap.MsgType_MSG_PINGPONG, pbmsgwrap.MsgSec_SEC_SIGN, wsmsg.GetMsgTrack(), &msg)
}

// null Msg - This is technically a "login" message, and is used to setup the presence object
func (d *Data) nullMsg(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_NONE {
		return errors.New("required security level not set")
	}

	modelPresence := presence.New()

	if err := modelPresence.Fetch(d.cbPool, wsmsg.GetMsgKey(), false); err != nil {
		return errors.New("get model presence failed")
	}

	signTest := security.CalcPBhash(modelPresence.GetSign(), wsmsg.GetMsgKey(), wsmsg.GetMsgTrack(), int(wsmsg.GetMsgType()), int(wsmsg.GetMsgSec()), wsmsg.GetMsgDate(), wsmsg.GetMsgData())
	if string(wsmsg.GetMsgHash()) != signTest {
		return errors.New("hash compare failed")
	}

	d.io.Lock()
	d.key = wsmsg.GetMsgKey()
	d.timeCheck = time.Now().Unix() + 150
	d.sign = modelPresence.GetSign()
	d.cipher = modelPresence.GetCipher()
	d.io.Unlock()

	modelPresence.SetServerid(d.serviceID)
	modelPresence.SetConnectionid(d.id)

	if err := modelPresence.Edit(d.cbPool); err != nil {
		return err
	}

	msg := []byte("")
	return d.Write(pbmsgwrap.MsgType_MSG_NULL, pbmsgwrap.MsgSec_SEC_NONE, wsmsg.GetMsgTrack(), &msg)
}

// latency Msg
func (d *Data) latencyMsg(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_SIGN {
		return errors.New("required security level not set")
	}

	msg, latency, err := mLatency.PackServer(wsmsg.GetMsgData())
	if err != nil {
		return err
	}

	d.latency = latency

	return d.Write(pbmsgwrap.MsgType_MSG_LATENCY, pbmsgwrap.MsgSec_SEC_SIGN, wsmsg.GetMsgTrack(), &msg)
}

// SrvTime Msg
func (d *Data) srvtimeMsg(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_SIGN {
		return errors.New("required security level not set")
	}

	msg, err := mSrvtime.Pack()
	if err != nil {
		return err
	}

	return d.Write(pbmsgwrap.MsgType_MSG_SRVTIME, pbmsgwrap.MsgSec_SEC_SIGN, wsmsg.GetMsgTrack(), &msg)
}

// asyncMsg places a message into message queue
func (d *Data) asyncMsg(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_SIGN && wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_CIPHER {
		return errors.New("required security level not set")
	}

	var queue string

	switch wsmsg.GetMsgType() {
	case pbmsgwrap.MsgType_MSG_CHAT_CTL:
		queue = "async_rpc_chat_ctl_queue"

	case pbmsgwrap.MsgType_MSG_CHAT_MSG:
		queue = "async_rpc_chat_msg_queue"

	case pbmsgwrap.MsgType_MSG_TELEMETRY:
		queue = "async_ff_telemetry_queue"

	default:
		d.prometheus.IncAsynctxDeliveryCount("async_unknown_queue", "bad")
		return errors.New("unknown async message type")
	}

	if err := d.asynctx.Schedule(d.id, d.key, wsmsg.GetMsgTrack(), queue, false, wsmsg.GetMsgData()); err != nil {
		d.prometheus.IncAsynctxDeliveryCount(queue, "bad")
		return err
	}
	d.prometheus.IncAsynctxDeliveryCount(queue, "good")

	return nil
}

// syncMsg pushes the message to the sync processing system
func (d *Data) syncMsg(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_SIGN || wsmsg.GetMsgSec() != pbmsgwrap.MsgSec_SEC_CIPHER {
		return errors.New("required security level not set")
	}

	d.log.Sugar().Debugf("syncMsg: %d", wsmsg.GetMsgType())

	return nil
}

// EOF
