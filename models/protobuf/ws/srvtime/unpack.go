package srvtime

import (
	"github.com/golang/protobuf/proto"
	pbsrvtime "github.com/zk101/nixie/proto/ws/srvtime"
)

// Unpack srvtime message
func Unpack(data []byte) (*pbsrvtime.SrvTime, error) {
	srvtime := pbsrvtime.SrvTime{}

	if err := proto.Unmarshal(data, &srvtime); err != nil {
		return nil, err
	}

	return &srvtime, nil
}

// EOF
