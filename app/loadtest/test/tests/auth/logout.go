package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/lib/httpd"
	"github.com/zk101/nixie/lib/security"
)

// Logout does a logout check
func Logout(local *data.Local) error {
	if local.Key == "" || local.Sign == "" || local.Cipher == "" {
		return errors.New("bad logout no signing information")
	}

	uri := local.AuthBase + "/auth/logout"

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

	if resp.StatusCode == 200 {
		local.Key = ""
		local.Sign = ""
		local.Cipher = ""
		return nil
	}

	return errors.New("bad logout status code")
}

// EOF
