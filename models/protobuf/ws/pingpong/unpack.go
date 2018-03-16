package pingpong

import (
	"github.com/golang/protobuf/proto"
	pbping "github.com/zk101/nixie/proto/ws/pingpong"
)

// Unpack retusn a PingPonig message pointer
func Unpack(data []byte) (*pbping.PingPong, error) {
	ping := pbping.PingPong{}

	if err := proto.Unmarshal(data, &ping); err != nil {
		return nil, err
	}

	return &ping, nil
}

// EOF
