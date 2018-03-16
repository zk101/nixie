package user

import (
	"errors"
	"time"

	"github.com/zk101/nixie/app/ws/prometheus"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Receive processes a mesage from a connection
func (d *Data) Receive() error {
	start := time.Now()
	wsmsg, msgLen, err := d.Read()
	if err != nil {
		d.conn.Close()
		return err
	}

	if err := d.checkPresence(wsmsg); err != nil {
		d.conn.Close()
		return err
	}

	switch wsmsg.GetMsgType() {
	case pbmsgwrap.MsgType_MSG_NULL:
		if err := d.nullMsg(wsmsg); err != nil {
			d.conn.Close()
			return err
		}

	case pbmsgwrap.MsgType_MSG_PINGPONG:
		if err := d.pingpingMsg(wsmsg); err != nil {
			d.conn.Close()
			return err
		}

	case pbmsgwrap.MsgType_MSG_LATENCY:
		if err := d.latencyMsg(wsmsg); err != nil {
			d.conn.Close()
			return err
		}

	case pbmsgwrap.MsgType_MSG_SRVTIME:
		if err := d.srvtimeMsg(wsmsg); err != nil {
			d.conn.Close()
			return err
		}

	case pbmsgwrap.MsgType_MSG_CHAT_CTL, pbmsgwrap.MsgType_MSG_CHAT_MSG:
		if err := d.asyncMsg(wsmsg); err != nil {
			d.conn.Close()
			return err
		}

	case pbmsgwrap.MsgType_MSG_TELEMETRY:
		if err := d.asyncMsg(wsmsg); err != nil {
			d.conn.Close()
			return err
		}

	default:
		d.conn.Close()
		return errors.New("unknown message type")
	}

	d.prometheus.IncReqCount(pbmsgwrap.MsgType_name[int32(wsmsg.GetMsgType())], prometheus.DirectionRX, pbmsgwrap.MsgSec_name[int32(wsmsg.GetMsgSec())])
	d.prometheus.ObserveReqSize(msgLen, pbmsgwrap.MsgType_name[int32(wsmsg.GetMsgType())], prometheus.DirectionRX, pbmsgwrap.MsgSec_name[int32(wsmsg.GetMsgSec())])
	d.prometheus.ObserveReqDuration(start, pbmsgwrap.MsgType_name[int32(wsmsg.GetMsgType())], prometheus.DirectionRX, pbmsgwrap.MsgSec_name[int32(wsmsg.GetMsgSec())])

	return nil
}

// EOF
