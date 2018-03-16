package presence

import (
	"errors"

	"github.com/zk101/nixie/lib/storage"
)

// Edit updates a Presence object
func (m *Model) Edit(cb storage.Couchbase) error {
	if m.new == nil {
		return nil
	}

	if m.cur.GetKey() != m.new.GetKey() {
		return errors.New("model presence key can not be modified")
	}

	if m.cur.GetUserid() != m.new.GetUserid() {
		return errors.New("model presence userid can not be modified")
	}

	if m.cur.GetDn() != m.new.GetDn() {
		return errors.New("model presence dn can not be modified")
	}

	data, err := m.Pack()
	if err != nil {
		return err
	}

	if _, err := cb.Replace(m.cur.Key, &data, m.curCas, 300); err != nil {
		return err
	}

	if _, err := cb.Touch(m.cur.Userid, 0, 300); err != nil {
		return err
	}

	if _, err := cb.Touch(m.cur.Dn, 0, 300); err != nil {
		return err
	}

	return nil
}

// EOF
