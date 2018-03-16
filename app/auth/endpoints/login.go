package endpoints

import (
	"io/ioutil"
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
	"github.com/zk101/nixie/models/couchbase/auth/presence"
	"github.com/zk101/nixie/models/ldap/auth/useredit"
	mLogin "github.com/zk101/nixie/models/protobuf/auth/login"
)

// Login checks a users creds and logs them in
func (c *Client) Login(r *http.Request, b *body.Data) {
	if r.Method != "PUT" {
		c.log.Sugar().Warnw("method not implemented", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
		b.Status = 501
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.log.Sugar().Warnw("read body failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	reqData, err := mLogin.UnpackRequest(&reqBody)
	if err != nil {
		c.log.Sugar().Warnw("unpack request failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	user := useredit.New()
	user.Base = c.ldapPool.GetBase()
	user.Username = reqData.GetUsername()

	if err := user.Fetch(c.ldapPool); err != nil {
		c.log.Sugar().Warnw("fetch useredit model failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	if err := c.ldapPool.Bind(user.DN, reqData.GetPassword()); err != nil {
		c.log.Sugar().Errorw("user authentication failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	modelPresence := presence.New()
	modelPresence.SetDn(user.DN)
	modelPresence.SetUserid(reqData.GetUsername())

	if err := modelPresence.Create(c.cbPool); err != nil {
		c.log.Sugar().Warnw("create presence failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Body = []byte("")
		b.Status = 401
		return
	}

	b.Body, err = mLogin.PackReply(modelPresence.GetKey(), modelPresence.GetSign(), modelPresence.GetCipher())
	if err != nil {
		c.log.Sugar().Warnw("pack reply failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Body = []byte("")
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
