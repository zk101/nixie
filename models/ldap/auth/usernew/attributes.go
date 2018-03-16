package usernew

import (
	"errors"
	"net"
	"net/mail"
	"strings"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/security"
	"github.com/zk101/nixie/lib/validation"
)

// userAttributes processes a model into ldap.EntryAttributes and returns it
func (m *Model) userAttributes() (*[]ldap.EntryAttribute, error) {
	if validation.CheckString(m.Name, `^[^!"#$%&()*+/:;<=>?@[\]^\x60{|}~[:cntrl:]\t\n\v\f\r]+$`) != true {
		return nil, errors.New("name validation failed")
	}

	if validation.CheckRuneCount(m.Name, 1, 128) != true {
		return nil, errors.New("name rune count failed")
	}

	if validation.CheckString(m.Username, "^[^[:punct:][:cntrl:][:space:]]+$") != true {
		return nil, errors.New("user id validation failed")
	}

	if validation.CheckRuneCount(m.Username, 1, 64) != true {
		return nil, errors.New("user id rune count failed")
	}

	if m.Password == "" {
		return nil, errors.New("password validation failed")
	}

	password, err := security.LDAPhashSSHA512(m.Password)
	if err != nil {
		return nil, err
	}

	attributes := []ldap.EntryAttribute{
		ldap.EntryAttribute{
			Name: "objectclass",
			Values: []string{
				"inetOrgPerson", "organizationalPerson", "person", "top",
			},
		},
		ldap.EntryAttribute{
			Name: "cn",
			Values: []string{
				m.UUID,
			},
		},
		ldap.EntryAttribute{
			Name: "sn",
			Values: []string{
				m.UUID,
			},
		},
		ldap.EntryAttribute{
			Name: "displayName",
			Values: []string{
				m.Name,
			},
		},
		ldap.EntryAttribute{
			Name: "uid",
			Values: []string{
				m.Username,
			},
		},
		ldap.EntryAttribute{
			Name: "userPassword",
			Values: []string{
				password,
			},
		},
	}

	if m.Email != "" {
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
	}

	if m.Mobile != "" {
		mobile := strings.Replace(m.Mobile, " +-()", "", -1)

		if validation.CheckString(mobile, "^[0-9]{1,15}$") != true {
			return nil, errors.New("mobile validation failed")
		}

		attributes = append(attributes, ldap.EntryAttribute{Name: "mobile", Values: []string{m.Mobile}})
	}

	if m.Address != "" {
		if validation.CheckString(m.Address, `^[^!"#$%&'()*+,.:;<=>?@[\]^_\x60{|}~\x00-\x09\x0B-\x1F\x7F\t\v\f\r]+$`) != true {
			return nil, errors.New("address validation failed")
		}

		if validation.CheckRuneCount(m.Address, 1, 256) != true {
			return nil, errors.New("address rune count failed")
		}

		attributes = append(attributes, ldap.EntryAttribute{Name: "postalAddress", Values: []string{m.Address}})
	}

	return &attributes, nil
}

// EOF
