package auth

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/lib/security"
	mLogin "github.com/zk101/nixie/models/protobuf/auth/login"
	pblogin "github.com/zk101/nixie/proto/auth/login"
)

// Login sends a login request and outputs the results
func Login(local *data.Local) error {
	if local.User == "" {
		return errors.New("login requires a user")
	}

	if local.Pass == "" {
		return errors.New("login requires a password")
	}

	uri := local.AuthBase + "/auth/login"

	loginBody, err := mLogin.PackRequest(local.User, local.Pass)
	if err != nil {
		return err
	}

	request, err := NewRequest("PUT", uri, bytes.NewReader(loginBody))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("bad login status code")
	}

	repBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	repData, err := mLogin.UnpackReply(&repBody)
	if err != nil {
		return err
	}

	if repData.GetError() == pblogin.LoginReply_FAIL {
		return errors.New("auth reply marked as fail")
	}

	if len(resp.Header["Authorization"]) != 1 {
		return errors.New("bad authorization header missing or more than one")
	}

	authBits := strings.Split(resp.Header["Authorization"][0], " ")
	if len(authBits) != 2 {
		return errors.New("bad authorization header does not consist of two parts")
	}

	if len(resp.Header["Auth-Date"]) != 1 {
		return errors.New("bad date header missing or more than one")
	}

	if authBits[1] != security.CalcAuthHash(repData.GetSign(), resp.Header["Auth-Date"][0], uri, repBody) {
		return errors.New("authorization hash check failed")
	}

	local.Key = repData.GetKey()
	local.Sign = repData.GetSign()
	local.Cipher = repData.GetCipher()

	return nil
}

// EOF
