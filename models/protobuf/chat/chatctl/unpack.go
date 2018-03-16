package chatctl

import (
	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/proto/chat"
)

// Unpack unmarshals a protobuf ChatCtl msg
func (m *Model) Unpack(data []byte) error {
	ctlData := chat.ChatCtl{}

	if err := proto.Unmarshal(data, &ctlData); err != nil {
		return err
	}

	*m = Model(ctlData)

	return nil
}

// EOF
