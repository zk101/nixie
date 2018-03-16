package storage

import "github.com/mavricknz/ldap"

// LDAP provides an interface for LDAP Server
type LDAP interface {
	Add(req *ldap.AddRequest) error
	Bind(username, password string) error
	Delete(delReq *ldap.DeleteRequest) error
	ModDn(req *ldap.ModDnRequest) error
	Modify(modReq *ldap.ModifyRequest) error
	Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error)
}

// EOF
