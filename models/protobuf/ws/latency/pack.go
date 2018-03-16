package latency

import (
	"errors"
	"time"

	"github.com/golang/protobuf/proto"
	pblatency "github.com/zk101/nixie/proto/ws/latency"
)

// PackClient returns
func PackClient() ([]byte, error) {
	latency := pblatency.Latency{
		ClientTime: time.Now().UnixNano(),
	}

	data, err := proto.Marshal(&latency)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// PackServer takes a current Latency and returns a packed server version
func PackServer(data []byte) ([]byte, int64, error) {
	latency, err := Unpack(data)
	if err != nil {
		return nil, 0, err
	}

	latency.ServerTime = time.Now().UnixNano()

	timeDiff := latency.ServerTime - latency.ClientTime
	if timeDiff < 0 {
		return nil, 0, errors.New("latency differance less than 0")
	}

	msg, err := proto.Marshal(latency)
	if err != nil {
		return nil, 0, err
	}

	return msg, timeDiff, nil
}

// EOF
