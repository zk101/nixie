package usernew

import (
	"errors"
	"fmt"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
)

// Create adds a new user to storage
func (m *Model) Create(storeLDAP storage.LDAP) error {
	//if m.DN == "" {
	//	return errors.New("model requires dn to be set")
	//}

	if m.Base == "" {
		return errors.New("model requires base to be set")
	}

	if m.Base == "" {
		return errors.New("model requires base to be set")
	}

	if m.Username == "" {
		return errors.New("model requires username to be set")
	}

	sr := ldap.NewSimpleSearchRequest(m.Base, ldap.ScopeWholeSubtree, fmt.Sprintf("(uid=%s)", m.Username), []string{"uid"})
	result, err := storeLDAP.Search(sr)
	if err != nil {
		return err
	}

	if len(result.Entries) != 0 {
		return errors.New("user exists")
	}

	addAttrs, err := m.userAttributes()
	if err != nil {
		return err
	}

	// TODO: This should be moved to the calling function.
	m.DN = "cn=" + m.UUID + ",ou=users," + m.Base

	req := ldap.NewAddRequest(m.DN)
	for _, attr := range *addAttrs {
		req.AddAttribute(&attr)
	}

	return storeLDAP.Add(req)
}

// EOF
