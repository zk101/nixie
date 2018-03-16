package auth

import (
	"errors"
	"net/http"

	"github.com/zk101/nixie/app/loadtest/data"
)

// Health does a health check
func Health(local *data.Local) error {
	uri := local.AuthBase + "/auth/health"

	request, err := NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("bad health status code")
	}

	return nil
}

// EOF
