package worker

import (
	"errors"
	"fmt"

	"github.com/mavricknz/ldap"
	"github.com/zk101/nixie/lib/validation"
	"github.com/zk101/nixie/models/couchbase/auth/presence"
	"github.com/zk101/nixie/models/ldap/chat/userchat"
	"github.com/zk101/nixie/models/protobuf/chat/chatctl"
)

// ctlNull deals with Contral Status messages
func (c *Client) ctlNull(msg *chatctl.Model) error {
	modelPresence := presence.New()
	if err := modelPresence.Fetch(c.cbPool, msg.Userid, true); err != nil {
		return errors.New("get model presence failed")
	}

	user := userchat.New()
	user.DN = modelPresence.GetDn()

	if err := user.Fetch(c.ldapPool); err != nil {
		return err
	}

	modelPresence.SetChatfriends(user.Friends)

	return modelPresence.Edit(c.cbPool)
}

// ctlSearch deals with Search msgs
func (c *Client) ctlSearch(msg *chatctl.Model) error {
	if validation.CheckString(string(msg.Userid), "^[^[:punct:][:cntrl:][:space:]]+$") != true {
		return errors.New("user id validation failed")
	}

	sr := ldap.NewSimpleSearchRequest(c.ldapPool.GetBase(), ldap.ScopeWholeSubtree, fmt.Sprintf("(uid=%s)", string(msg.Userid)), []string{"uid"})
	result, err := c.ldapPool.Search(sr)
	if err == nil {
		return err
	}

	if len(result.Entries) != 1 {
		return errors.New("chat search failed user does not exist")
	}

	sr = ldap.NewSimpleSearchRequest(fmt.Sprintf("cn=friends,ou=chat,%s", result.Entries[0].DN), ldap.ScopeBaseObject, "(objectClass=groupOfNames)", []string{"cn"})
	result, err = c.ldapPool.Search(sr)
	if err == nil {
		return err
	}

	if len(result.Entries) != 1 {
		return errors.New("chat search failed user not registered for chat")
	}

	return nil
}

// EOF
