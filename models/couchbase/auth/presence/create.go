package presence

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/zk101/nixie/lib/storage"
)

// Create sets up a Presence PB, inserts it and returns a pointer to the new Presence
func (m *Model) Create(cb storage.Couchbase) error {
	if m.new == nil {
		return errors.New("model presence create on empty model")
	}

	if m.new.GetKey() == "" {
		return errors.New("model presence key must be set")
	}

	if m.new.GetUserid() == "" {
		return errors.New("model presence userid must be set")
	}

	if m.new.GetDn() == "" {
		return errors.New("model presence dn must be set")
	}

	data, err := m.Pack()
	if err != nil {
		return err
	}

	if _, err := cb.Insert(m.new.GetKey(), &data, 300); err != nil {
		return err
	}

	presenceID := newPBpresenceID(m.new.Key)

	data, err = proto.Marshal(presenceID)
	if err != nil {
		return err
	}

	if _, err := cb.Upsert(m.new.GetUserid(), &data, 300); err != nil {
		return err
	}

	if _, err := cb.Upsert(m.new.GetDn(), &data, 300); err != nil {
		return err
	}

	return nil
}

// EOF
