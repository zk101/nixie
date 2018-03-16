package endpoints

import (
	"io/ioutil"
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
	"github.com/zk101/nixie/models/ldap/auth/usernew"
	mRegister "github.com/zk101/nixie/models/protobuf/auth/register"
)

// Register adds a user to the database
func (c *Client) Register(r *http.Request, b *body.Data) {
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

	reqData, err := mRegister.Unpack(&reqBody)
	if err != nil {
		c.log.Sugar().Warnw("unmarshal body failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	modelUser := usernew.New()
	modelUser.Base = c.ldapPool.GetBase()
	modelUser.Name = reqData.GetName()
	modelUser.Username = reqData.GetUsername()
	modelUser.Password = reqData.GetPassword()
	modelUser.Email = reqData.GetEmail()
	modelUser.Mobile = reqData.GetMobile()
	modelUser.Address = reqData.GetAddress()

	if err := modelUser.Create(c.ldapPool); err != nil {
		c.log.Sugar().Warnw("add user failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		b.Status = 401
		return
	}

	b.Status = 200
}

// EOF
