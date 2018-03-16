package presence

import "github.com/zk101/nixie/lib/storage"

// Touch updates an expiry on a the Presence objects
func (m *Model) Touch(cb storage.Couchbase) error {
	if _, err := cb.Touch(m.cur.GetKey(), 0, 300); err != nil {
		return err
	}

	if _, err := cb.Touch(m.cur.GetUserid(), 0, 300); err != nil {
		return err
	}

	if _, err := cb.Touch(m.cur.GetDn(), 0, 300); err != nil {
		return err
	}

	return nil
}

// EOF
