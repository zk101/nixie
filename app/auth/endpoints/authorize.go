package endpoints

import (
	"errors"
	"net/http"

	"github.com/zk101/nixie/app/auth/body"
	"github.com/zk101/nixie/lib/security"
	"github.com/zk101/nixie/models/couchbase/auth/presence"
)

// authorizeMessage is a wrapper around a series of common message validation tasks
func (c *Client) authorizeMessage(r *http.Request, b *body.Data) (*presence.Model, error) {
	key, hash, date, err := security.ValidateAuthorization(r, c.authTimeout)
	if err != nil {
		return nil, errors.New("message validatation failed")
	}

	modelPresence := presence.New()
	if err := modelPresence.Fetch(c.cbPool, key, false); err != nil {
		return nil, errors.New("get model presence failed")
	}

	if hash != security.CalcAuthHash(modelPresence.GetSign(), date, b.URL, []byte("")) {
		return nil, errors.New("hash validatation failed")
	}

	return modelPresence, nil
}

// EOF
