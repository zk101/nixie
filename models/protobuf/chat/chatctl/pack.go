package chatctl

import (
	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/proto/chat"
)

// Pack prepares model for transport
func (m *Model) Pack() ([]byte, error) {
	chatctl := chat.ChatCtl(*m)

	return proto.Marshal(&chatctl)
}

// EOF
