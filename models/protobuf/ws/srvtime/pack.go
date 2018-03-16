package srvtime

import (
	"time"

	"github.com/golang/protobuf/proto"
	pbsrvtime "github.com/zk101/nixie/proto/ws/srvtime"
)

// Pack srvtime message
func Pack() ([]byte, error) {
	srvtime := &pbsrvtime.SrvTime{
		ServerTime: time.Now().Format("Mon Jan 02 15:04:05 -0700 MST 2006"),
	}

	msg, err := proto.Marshal(srvtime)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// EOF
