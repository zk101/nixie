package presence

import (
	"strings"
	"time"

	"github.com/couchbase/gocb"
	uuid "github.com/satori/go.uuid"
	pbPresence "github.com/zk101/nixie/proto/presence"
)

// Model is an exact copy of the protobuf generated Presence struct minus the annotations.  This allows conversion to just work.
type Model struct {
	cur    *pbPresence.Presence
	curCas gocb.Cas
	new    *pbPresence.Presence
}

// New returns a pointer to an empty presence Model
func New() *Model {
	return &Model{
		cur: newPBpresence(),
	}
}

// newPBpresence creates a new Presence Protobuf struct
func newPBpresence() *pbPresence.Presence {
	return &pbPresence.Presence{
		Date:        time.Now().Unix(),
		Key:         strings.Replace(uuid.NewV4().String(), "-", "", -1),
		Sign:        strings.Replace(uuid.NewV4().String()+uuid.NewV4().String(), "-", "", -1),
		Cipher:      strings.Replace(uuid.NewV4().String(), "-", "", -1),
		Chatfriends: make(map[string]bool, 0),
	}
}

// newPBpresenceID creates a new Presence Protobuf struct
func newPBpresenceID(key string) *pbPresence.PresenceID {
	return &pbPresence.PresenceID{
		Key: key,
	}
}

// EOF
