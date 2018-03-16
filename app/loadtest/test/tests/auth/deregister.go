package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/security"
)

// Deregister does a refresh check
func Deregister(local *data.Local) error {
	if local.Key == "" || local.Sign == "" || local.Cipher == "" {
		return errors.New("bad refresh no signing information")
	}

	uri := local.AuthBase + "/auth/deregister"

	request, err := NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	date := httpd.HTTPdate()
	request.Header.Add("Authorization", fmt.Sprintf("%s %s", local.Key, security.CalcAuthHash(local.Sign, date, uri, []byte(""))))
	request.Header.Add("Auth-Date", date)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("bad deregister status code")
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

	if authBits[1] != security.CalcAuthHash(local.Sign, resp.Header["Auth-Date"][0], uri, []byte("")) {
		return errors.New("authorization hash check failed")
	}

	local.Key = ""
	local.Sign = ""
	local.Cipher = ""

	return nil
}

// EOF
