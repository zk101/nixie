package endpoints

import (
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
	"github.com/zk101/nixie/models/ldap/auth/useredit"
)

// Deregister tests if a user is logged in
func (c *Client) Deregister(r *http.Request, b *body.Data) {
	if r.Method != "GET" {
		c.log.Sugar().Warnw("method not implemented", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
		b.Status = 501
		return
	}

	modelPresence, err := c.authorizeMessage(r, b)
	if err != nil {
		c.log.Sugar().Errorw("message authorization failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	user := useredit.New()
	user.DN = modelPresence.GetDn()

	if err := user.Remove(c.ldapPool); err != nil {
		c.log.Sugar().Warnw("fetch useredit model failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 500
		return
	}

	if err := modelPresence.Remove(c.cbPool); err != nil {
		c.log.Sugar().Warnw("remove model presence failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 500
		return
	}

	b.Status = 200
	b.Auth = true
	b.Key = modelPresence.GetKey()
	b.Sign = modelPresence.GetSign()
	b.Cipher = modelPresence.GetCipher()
}

// EOF
