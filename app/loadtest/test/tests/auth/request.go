package auth

import (
	"io"
	"net/http"
)

// NewRequest sets up a new http request and returns a poijnter to it
func NewRequest(method, uri string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	request.Close = true
	request.Header.Set("Accept", "application/octet-stream")
	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("User-Agent", "loadtest")

	return request, nil
}

// EOF
