package pingpong

import (
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	pbping "github.com/zk101/nixie/proto/ws/pingpong"
)

// PackClient returns a []byte
func PackClient() ([]byte, string, error) {
	ping := &pbping.PingPong{
		PingPong: pbping.PingPong_PING,
		Sequence: uuid.NewV4().String(),
	}

	data, err := proto.Marshal(ping)
	if err != nil {
		return nil, "", err
	}

	return data, ping.GetSequence(), nil
}

// PackServer returns a []byte
func PackServer(data []byte) ([]byte, error) {
	pong, err := Unpack(data)
	if err != nil {
		return nil, err
	}

	pong.PingPong = pbping.PingPong_PONG

	msg, err := proto.Marshal(pong)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// EOF
