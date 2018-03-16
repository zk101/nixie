package body

import (
	"net/http"

	"github.com/zk101/nixie/lib/httpd"
)

// Data is a storage structure to collect varies message pieces together
type Data struct {
	URL     string
	Status  int
	Date    string
	Content string
	Body    []byte
	Auth    bool
	Key     string
	Sign    string
	Cipher  string
}

// InitData returns an initialised body
func InitData(r *http.Request) *Data {
	b := &Data{httpd.URL(r), 404, httpd.HTTPdate(), "application/octet-stream", []byte(""), false, "", "", ""}

	return b
}

// EOF
