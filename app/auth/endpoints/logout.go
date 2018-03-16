package endpoints

import (
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
)

// Logout logs a user out
func (c *Client) Logout(r *http.Request, b *body.Data) {
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

	if err := modelPresence.Remove(c.cbPool); err != nil {
		c.log.Sugar().Warnw("remove model presence failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	b.Status = 200
}

// EOF
