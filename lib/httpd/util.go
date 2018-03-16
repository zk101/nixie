package httpd

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// HTTPdate returns a HTTP Date
func HTTPdate() string {
	t := time.Now()

	return t.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

// URL unmangles the varies http package url pieces into a single url.
func URL(r *http.Request) string {
	var proto string

	if len(r.Header["X-Forwarded-Proto"]) == 1 {
		proto = fmt.Sprintf("%s://", r.Header["X-Forwarded-Proto"][0])
	} else {
		if r.TLS == nil {
			proto = "http://"
		} else {
			proto = "https://"
		}
	}
	uri, _ := url.QueryUnescape(r.RequestURI)

	return proto + r.Host + uri
}

// EOF
