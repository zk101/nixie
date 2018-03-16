package chatctl

import (
	"time"

	"github.com/zk101/nixie/proto/chat"
)

// Model is an exact copy of the protobuf generated Presence struct minus the annotations.  This allows conversion to just work.
type Model struct {
	Type   chat.ChatCtlType
	Status chat.ChatCtlStatus
	Userid string
	Date   int64
}

// New returns a pointer to an empty presence Model
func New() *Model {
	return &Model{
		Date: time.Now().Unix(),
	}
}

// EOF
