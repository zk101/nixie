package chatmsg

import (
	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/proto/chat"
)

// Pack prepares model for transport
func (m *Model) Pack() ([]byte, error) {
	chatmsg := chat.ChatMsg(*m)

	return proto.Marshal(&chatmsg)
}

// EOF
