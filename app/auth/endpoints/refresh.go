package endpoints

import (
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
)

// Refresh tests if a user is logged in
func (c *Client) Refresh(r *http.Request, b *body.Data) {
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

	if err := modelPresence.Touch(c.cbPool); err != nil {
		c.log.Sugar().Warnw("touch model presence failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	b.Status = 200
	b.Auth = true
	b.Key = modelPresence.GetKey()
	b.Sign = modelPresence.GetSign()
	b.Cipher = modelPresence.GetCipher()
}

// EOF
