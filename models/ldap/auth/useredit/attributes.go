package useredit

import (
	"errors"
	"net"
	"net/mail"
	"strings"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/validation"
)

// userAttributes processes a model into ldap.EntryAttributes and returns it
func (m *Model) userAttributes(oldValues *Model) (*[]ldap.Mod, error) {
	var mod *ldap.Mod
	mods := make([]ldap.Mod, 0)

	if m.Name == "" {
		return nil, errors.New("name can not be empty")
	}

	if m.Name != oldValues.Name {
		if validation.CheckString(m.Name, `^[^!"#$%&()*+/:;<=>?@[\]^\x60{|}~[:cntrl:]\t\n\v\f\r]+$`) != true {
			return nil, errors.New("name validation failed")
		}

		if validation.CheckRuneCount(m.Name, 1, 128) != true {
			return nil, errors.New("name rune count failed")
		}

		mod = ldap.NewMod(ldap.ModReplace, "displayName", []string{m.Name})
		mods = append(mods, *mod)
	}

	if m.Email != oldValues.Email {
		if m.Email == "" {
			mod := ldap.NewMod(ldap.ModDelete, "mail", []string{})
			mods = append(mods, *mod)
		} else {
			e, err := mail.ParseAddress(m.Email)
			if err != nil {
				return nil, err
			}
			m.Email = e.Address

			addressParts := strings.Split(m.Email, "@")
			if len(addressParts) != 2 {
				return nil, errors.New("email address did not split into two parts")
			}

			if _, err := net.LookupMX(addressParts[1]); err != nil {
				return nil, err
			}

			if oldValues.Email == "" {
				mod = ldap.NewMod(ldap.ModAdd, "mail", []string{m.Email})
			} else {
				mod = ldap.NewMod(ldap.ModReplace, "mail", []string{m.Email})
			}
			mods = append(mods, *mod)
		}
	}

	if m.Mobile != oldValues.Mobile {
		if m.Mobile == "" {
			mod := ldap.NewMod(ldap.ModDelete, "mobile", []string{})
			mods = append(mods, *mod)
		} else {
			mobile := strings.Replace(m.Mobile, " +-()", "", -1)

			if validation.CheckString(mobile, "^[0-9]{1,15}$") != true {
				return nil, errors.New("mobile validation failed")
			}

			if oldValues.Mobile == "" {
				mod = ldap.NewMod(ldap.ModAdd, "mobile", []string{m.Mobile})
			} else {
				mod = ldap.NewMod(ldap.ModReplace, "mobile", []string{m.Mobile})
			}
			mods = append(mods, *mod)
		}
	}

	if m.Address != oldValues.Address {
		if m.Address == "" {
			mod = ldap.NewMod(ldap.ModDelete, "postalAddress", []string{})
			mods = append(mods, *mod)
		} else {
			if validation.CheckString(m.Address, `^[^!"#$%&'()*+,.:;<=>?@[\]^_\x60{|}~\x00-\x09\x0B-\x1F\x7F\t\v\f\r]+$`) != true {
				return nil, errors.New("address validation failed")
			}

			if validation.CheckRuneCount(m.Address, 1, 256) != true {
				return nil, errors.New("address rune count failed")
			}

			if oldValues.Address == "" {
				mod = ldap.NewMod(ldap.ModAdd, "postalAddress", []string{m.Address})
			} else {
				mod = ldap.NewMod(ldap.ModReplace, "postalAddress", []string{m.Address})
			}
			mods = append(mods, *mod)
		}
	}

	return &mods, nil
}

// EOF
