package usernew

import (
	"errors"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
)

// Remove deletes a user.  This function is primarily used for testing, and is exported to allow rolling back a successful Create
func (m *Model) Remove(storeLDAP storage.LDAP) error {
	if m.DN == "" {
		return errors.New("model requires dn to be set")
	}

	sr := ldap.NewSimpleSearchRequest(m.DN, ldap.ScopeWholeSubtree, "(objectClass=*)", []string{"none"})
	result, err := storeLDAP.Search(sr)
	if err != nil {
		return err
	}

	for i := len(result.Entries) - 1; i >= 0; i-- {
		if err := storeLDAP.Delete(ldap.NewDeleteRequest(result.Entries[i].DN)); err != nil {
			return err
		}
	}

	return nil
}

// EOF
