package chatmsg

import (
	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/proto/chat"
)

// Unpack unmarshals a protobuf ChatMsg msg
func (m *Model) Unpack(data []byte) error {
	msgData := chat.ChatMsg{}

	if err := proto.Unmarshal(data, &msgData); err != nil {
		return err
	}

	*m = Model(msgData)

	return nil
}

// EOF
