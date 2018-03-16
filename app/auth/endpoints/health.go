package endpoints

import (
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
	"github.com/zk101/nixie/lib/couchbase"
	"github.com/zk101/nixie/lib/ldap"
	"github.com/zk101/nixie/lib/network"
)

// Health tests the server is working
func (c *Client) Health(r *http.Request, b *body.Data) {
	if err := network.CheckHTTPhost(r.RemoteAddr, `^((127\.0\.0\.)|(10\.)|(192\.168\.)|(172\.1[6-9]\.)|(172\.2[0-9]\.)|(172\.3[0-1]\.))`); err != nil {
		c.log.Sugar().Errorw("remote address test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		return
	}

	if r.Method != "GET" {
		c.log.Sugar().Warnw("method not implemented", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
		b.Status = 501
		return
	}

	testCB := couchbase.NewClient(c.cbPool.GetBucketConfig())
	if err := testCB.Test(); err != nil {
		c.log.Sugar().Errorw("couchbase test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 503
		return
	}

	testLDAPro := ldap.NewClient(c.ldapPool.GetROconfig(), c.ldapPool.GetCAcertPool())
	if err := testLDAPro.Test(); err != nil {
		c.log.Sugar().Errorw("ldap readonly test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 503
		return
	}

	testLDAPrw := ldap.NewClient(c.ldapPool.GetRWconfig(), c.ldapPool.GetCAcertPool())
	if err := testLDAPrw.Test(); err != nil {
		c.log.Sugar().Errorw("ldap readwrite test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 503
		return
	}

	b.Status = 200
}

// EOF
