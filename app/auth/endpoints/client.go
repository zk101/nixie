package endpoints

import (
	"errors"

	cbpool "github.com/zk101/nixie/lib/couchbase/pool"
	ldappool "github.com/zk101/nixie/lib/ldap/pool"
	"github.com/zk101/nixie/lib/logging"
)

// Client holds pointers to common application resources used by the endpoints
type Client struct {
	authTimeout uint32
	cbPool      *cbpool.Client
	ldapPool    *ldappool.Client
	log         *logging.Client
}

// NewClient returns a configured Client struct
func NewClient(authTimeout uint32, cbPool *cbpool.Client, ldapPool *ldappool.Client, log *logging.Client) (*Client, error) {
	if authTimeout < 1 {
		return nil, errors.New("authtime must be atleast 1")
	}

	if cbPool == nil {
		return nil, errors.New("require couchbase pool")
	}

	if ldapPool == nil {
		return nil, errors.New("require ldap pool")
	}

	if log == nil {
		return nil, errors.New("require a configured logger")
	}

	return &Client{
		authTimeout: authTimeout,
		cbPool:      cbPool,
		ldapPool:    ldapPool,
		log:         log,
	}, nil
}

// EOF
