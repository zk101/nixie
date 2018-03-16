package endpoints

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zk101/nixie/lib/network"
)

// Metrics wraps the Prometheus handler
func (c *Client) Metrics(w http.ResponseWriter, r *http.Request) bool {
	if err := network.CheckHTTPhost(r.RemoteAddr, `^((127\.0\.0\.)|(10\.)|(192\.168\.)|(172\.1[6-9]\.)|(172\.2[0-9]\.)|(172\.3[0-1]\.))`); err != nil {
		c.log.Sugar().Errorw("remote address test failed", "path", r.URL.Path, "error", err.Error(), "remote_addr", r.RemoteAddr)
		return false
	}

	c.log.Sugar().Debugw("request processed", "status", 200, "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)

	h := promhttp.Handler()
	h.ServeHTTP(w, r)

	return true
}

// EOF
