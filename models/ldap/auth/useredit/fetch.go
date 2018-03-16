package useredit

import (
	"errors"
	"fmt"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
)

// Fetch gets a User object
func (m *Model) Fetch(storeLDAP storage.LDAP) error {
	if m.Base == "" {
		return errors.New("model requires base to be set")
	}

	if m.Username == "" {
		return errors.New("model requires username to be set")
	}

	sr := ldap.NewSimpleSearchRequest(m.Base, ldap.ScopeWholeSubtree, fmt.Sprintf("(uid=%s)", m.Username), []string{"displayName", "mail", "mobile", "postalAddress"})
	result, err := storeLDAP.Search(sr)
	if err != nil {
		return err
	}

	if len(result.Entries) != 1 {
		return errors.New("user lookup failed")
	}

	m.DN = result.Entries[0].DN
	m.Name = result.Entries[0].GetAttributeValue("displayName")
	m.Email = result.Entries[0].GetAttributeValue("mail")
	m.Mobile = result.Entries[0].GetAttributeValue("mobile")
	m.Address = result.Entries[0].GetAttributeValue("postalAddress")

	return nil
}

// EOF
