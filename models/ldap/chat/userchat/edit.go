package userchat

import (
	"errors"
	"fmt"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/storage"
)

// Edit updates a User object
func (m *Model) Edit(storeLDAP storage.LDAP) error {
	if m.DN == "" {
		return errors.New("model requires dn to be set")
	}

	data := make([]string, 0)
	for key := range m.Friends {
		data = append(data, key)
	}

	mod := ldap.NewMod(ldap.ModReplace, "member", data)
	req := ldap.NewModifyRequest(fmt.Sprintf("cn=friends,ou=chat,%s", m.DN))
	mods := make([]ldap.Mod, 0)
	mods = append(mods, *mod)

	req.AddMods(mods)

	return storeLDAP.Modify(req)
}

// EOF
