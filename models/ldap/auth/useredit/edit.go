package useredit

import (
	"errors"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
)

// Edit updates a User object
func (m *Model) Edit(storeLDAP storage.LDAP) error {
	if m.Base == "" {
		return errors.New("model requires base to be set")
	}

	if m.Username == "" {
		return errors.New("model requires username to be set")
	}

	oldValues := New()
	oldValues.Base = m.Base
	oldValues.Username = m.Username

	if err := oldValues.Fetch(storeLDAP); err != nil {
		return err
	}

	req := ldap.NewModifyRequest(m.DN)
	mods, err := m.userAttributes(oldValues)
	if err != nil {
		return err
	}

	req.AddMods(*mods)

	return storeLDAP.Modify(req)
}

// EOF
