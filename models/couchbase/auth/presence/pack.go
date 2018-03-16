package presence

import (
	"github.com/golang/protobuf/proto"
)

// Pack returns a []byte packed protobuf or an error
func (m *Model) Pack() ([]byte, error) {
	return proto.Marshal(m.new)
}

// EOF
