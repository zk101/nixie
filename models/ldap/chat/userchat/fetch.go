package userchat

import (
	"errors"
	"fmt"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
	"github.com/zk101/nixie/lib/utilities"
)

// Fetch gets a User object
func (m *Model) Fetch(storeLDAP storage.LDAP) error {
	if m.DN == "" {
		return errors.New("model requires dn to be set")
	}

	sr := ldap.NewSimpleSearchRequest(fmt.Sprintf("cn=friends,ou=chat,%s", m.DN), ldap.ScopeBaseObject, "(objectClass=groupOfNames)", []string{"member"})
	result, err := storeLDAP.Search(sr)

	if err != nil {
		retCode := utilities.ExtractLDAPerrorCode(err)
		switch retCode {
		case 32:
			if err := m.Create(storeLDAP); err != nil {
				return err
			}
			m.Friends[m.DN] = true

		default:
			return err
		}
	} else {
		for _, value := range result.Entries {
			if value.GetAttributeIndex("member") < 0 {
				continue
			}

			for _, member := range value.GetAttributeValues("member") {
				if member == m.DN {
					m.Friends[m.DN] = true
					continue
				}

				sr2 := ldap.NewSimpleSearchRequest(fmt.Sprintf("ou=friends,ou=chat,%s", member), ldap.ScopeBaseObject, fmt.Sprintf("(member=%s)", member), []string{"member"})
				result2, err := storeLDAP.Search(sr2)
				if err != nil {
					m.Friends[member] = false
					continue
				}

				if len(result2.Entries) == 1 {
					m.Friends[member] = true
					continue
				}

				m.Friends[member] = false
			}
		}
	}

	return nil
}

// EOF
