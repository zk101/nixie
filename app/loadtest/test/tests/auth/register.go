package auth

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/zk101/nixie/app/loadtest/data"
	mRegister "github.com/zk101/nixie/models/protobuf/auth/register"
)

// Register does a refresh check
func Register(local *data.Local) error {
	if local.User == "" {
		return errors.New("register requires a user")
	}

	if local.Pass == "" {
		return errors.New("register requires a password")
	}

	if local.Name == "" {
		return errors.New("register requires a name")
	}

	uri := local.AuthBase + "/auth/register"

	registerBody, err := mRegister.Pack(local.User, local.Pass, local.Name)
	if err != nil {
		return err
	}

	request, err := NewRequest("PUT", uri, bytes.NewReader(registerBody))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("bad register status code")
	}

	return nil
}

// EOF
