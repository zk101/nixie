package presence

import (
	"github.com/golang/protobuf/proto"
)

// Unpack returns a []byte packed protobuf or an error
func (m *Model) Unpack(data []byte) error {
	if err := proto.Unmarshal(data, m.cur); err != nil {
		return err
	}

	if m.cur.GetChatfriends() == nil {
		m.cur.Chatfriends = make(map[string]bool, 0)
	}

	return nil
}

// EOF
