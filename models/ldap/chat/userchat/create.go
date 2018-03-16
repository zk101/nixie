package userchat

import (
	"errors"
	"fmt"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
	"github.com/zk101/nixie/lib/utilities"
)

// Create sets up the ldap entries for chat
func (m *Model) Create(storeLDAP storage.LDAP) error {
	if m.DN == "" {
		return errors.New("model requires dn to be set")
	}

	sr := ldap.NewSimpleSearchRequest(fmt.Sprintf("ou=chat,%s", m.DN), ldap.ScopeBaseObject, "(objectClass=groupOfNames)", []string{"member"})
	_, err := storeLDAP.Search(sr)
	if err != nil {
		retCode := utilities.ExtractLDAPerrorCode(err)

		switch retCode {
		case ldap.LDAPResultNoSuchObject:
			attributes := []ldap.EntryAttribute{
				ldap.EntryAttribute{
					Name: "objectclass",
					Values: []string{
						"organizationalUnit", "top",
					},
				},
				ldap.EntryAttribute{
					Name: "ou",
					Values: []string{
						"chat",
					},
				},
			}

			req := ldap.NewAddRequest(fmt.Sprintf("ou=chat,%s", m.DN))
			for _, attr := range attributes {
				req.AddAttribute(&attr)
			}

			if err := storeLDAP.Add(req); err != nil {
				return err
			}

		default:
			return err
		}
	}

	sr = ldap.NewSimpleSearchRequest(fmt.Sprintf("cn=friends,ou=chat,%s", m.DN), ldap.ScopeBaseObject, "(objectClass=groupOfNames)", []string{"member"})
	_, err = storeLDAP.Search(sr)
	if err != nil {
		retCode := utilities.ExtractLDAPerrorCode(err)

		switch retCode {
		case ldap.LDAPResultNoSuchObject:
			attributes := []ldap.EntryAttribute{
				ldap.EntryAttribute{
					Name: "objectclass",
					Values: []string{
						"groupOfNames", "top",
					},
				},
				ldap.EntryAttribute{
					Name: "cn",
					Values: []string{
						"friends",
					},
				},
				ldap.EntryAttribute{
					Name: "member",
					Values: []string{
						m.DN,
					},
				},
			}

			req := ldap.NewAddRequest(fmt.Sprintf("cn=friends,ou=chat,%s", m.DN))
			for _, attr := range attributes {
				req.AddAttribute(&attr)
			}

			if err := storeLDAP.Add(req); err != nil {
				return err
			}

		default:
			return err
		}
	}

	return nil
}

// EOF
