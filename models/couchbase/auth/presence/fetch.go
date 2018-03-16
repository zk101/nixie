package presence

import (
	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/lib/storage"
)

// Fetch retrieves a presence object, where isIdent switches between a m.Key lookup and a Userid or DN lookup
func (m *Model) Fetch(cb storage.Couchbase, key string, isIdent bool) error {
	var (
		data []byte
		err  error
	)

	if m.curCas, err = cb.Get(key, &data); err != nil {
		return err
	}

	if isIdent {
		presenceID := newPBpresenceID("")

		if err := proto.Unmarshal(data, presenceID); err != nil {
			return err
		}

		if m.curCas, err = cb.Get(presenceID.GetKey(), &data); err != nil {
			return err
		}
	}

	return m.Unpack(data)
}

// EOF
