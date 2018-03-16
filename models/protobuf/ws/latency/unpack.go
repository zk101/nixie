package latency

import (
	"github.com/golang/protobuf/proto"
	pblatency "github.com/zk101/nixie/proto/ws/latency"
)

// Unpack latency message
func Unpack(data []byte) (*pblatency.Latency, error) {
	latency := pblatency.Latency{}

	if err := proto.Unmarshal(data, &latency); err != nil {
		return nil, err
	}

	return &latency, nil
}

// EOF
