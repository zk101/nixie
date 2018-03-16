package user

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// Read checks and unpacks the buffer
func (d *Data) Read() (*pbmsgwrap.MsgWrap, int64, error) {
	buf, opCode, msgLen, err := d.read()
	if err != nil {
		return nil, 0, err
	}

	var msgwarp *pbmsgwrap.MsgWrap

	switch opCode {
	case 0x2:
		msgwarp, err = msgwrap.Unpack(buf, d.sign, d.cipher)
		if err != nil {
			return nil, msgLen, err
		}
	case 0x8:
		closeMsg := make([]byte, 2)
		binary.BigEndian.PutUint16(closeMsg, 1000)

		if err := d.write(&closeMsg, ws.OpClose); err != nil {
			return nil, msgLen, err
		}
		return nil, msgLen, errors.New("websocket client requested close")
	default:
		return nil, msgLen, errors.New("websocket opcode not supported")
	}

	return msgwarp, msgLen, nil
}

// read pulls data off the wire
func (d *Data) read() (*[]byte, ws.OpCode, int64, error) {
	d.io.Lock()
	defer d.io.Unlock()

	h, r, err := wsutil.NextReader(d.conn, ws.StateServerSide)
	if err != nil {
		return nil, h.OpCode, 0, err
	}

	buf := make([]byte, h.Length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, h.OpCode, 0, err
	}

	return &buf, h.OpCode, h.Length, nil
}

// EOF
