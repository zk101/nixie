package presence

import "github.com/zk101/nixie/lib/storage"

// Remove deletes a Presence objects
func (m *Model) Remove(cb storage.Couchbase) error {
	if _, err := cb.Remove(m.cur.GetKey(), 0); err != nil {
		return err
	}

	if _, err := cb.Remove(m.cur.GetUserid(), 0); err != nil {
		return err
	}

	if _, err := cb.Remove(m.cur.GetDn(), 0); err != nil {
		return err
	}

	return nil
}

// EOF
