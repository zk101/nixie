package user

import (
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/zk101/nixie/app/ws/prometheus"
	"github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Write takes a *[]byte and packs it for sending.
func (d *Data) Write(msgType pbmsgwrap.MsgType, msgSec pbmsgwrap.MsgSec, msgTrack string, msgData *[]byte) error {
	start := time.Now()

	msgwarp, err := msgwrap.Pack(msgType, msgSec, d.key, d.sign, d.cipher, msgTrack, msgData)
	if err != nil {
		return err
	}

	if err := d.write(msgwarp, ws.OpBinary); err != nil {
		return err
	}

	d.prometheus.IncReqCount(pbmsgwrap.MsgType_name[int32(msgType)], prometheus.DirectionTX, pbmsgwrap.MsgSec_name[int32(msgSec)])
	d.prometheus.ObserveReqSize(int64(len(*msgwarp)), pbmsgwrap.MsgType_name[int32(msgType)], prometheus.DirectionTX, pbmsgwrap.MsgSec_name[int32(msgSec)])
	d.prometheus.ObserveReqDuration(start, pbmsgwrap.MsgType_name[int32(msgType)], prometheus.DirectionTX, pbmsgwrap.MsgSec_name[int32(msgSec)])

	return nil
}

// write does the actual socket writing
func (d *Data) write(data *[]byte, opCode ws.OpCode) error {
	w := wsutil.NewWriter(d.conn, ws.StateServerSide, opCode)

	d.io.Lock()
	defer d.io.Unlock()

	if _, err := w.Write(*data); err != nil {
		return err
	}

	return w.Flush()
}

// EOF
