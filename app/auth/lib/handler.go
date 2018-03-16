package lib

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zk101/nixie/app/auth/body"
	"github.com/zk101/nixie/lib/security"
)

// HandlerBase deals with incoming requests
func (c *Core) HandlerBase(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	b := body.InitData(r)

	if err := security.ValidateHeaders(r); err != nil {
		c.Clients.Log.Sugar().Errorw("bad request headers", "path", r.URL.Path, "method", r.Method, "error", err.Error(), "remote_addr", r.RemoteAddr)
	} else {
		switch r.URL.Path {
		case "/auth/v1/deregister", "/auth/deregister":
			c.Endpoints.Deregister(r, b)

		case "/auth/v1/health", "/auth/health":
			c.Endpoints.Health(r, b)

		case "/auth/v1/login", "/auth/login":
			c.Endpoints.Login(r, b)

		case "/auth/v1/logout", "/auth/logout":
			c.Endpoints.Logout(r, b)

		case "/auth/v1/refresh", "/auth/refresh":
			c.Endpoints.Refresh(r, b)

		case "/auth/v1/register", "/auth/register":
			c.Endpoints.Register(r, b)

		case "/metrics":
			if okay := c.Endpoints.Metrics(w, r); okay == true {
				c.Clients.Prometheus.IncReqCount(200, r.Method, r.URL.Path)
				c.Clients.Prometheus.ObserveReqDuration(&start, r.URL.Path)
				return
			}

		default:
			c.Clients.Log.Sugar().Errorw("bad path", "path", r.URL.Path, "method", r.Method, "remote_addr", r.RemoteAddr)
		}
	}

	if b.Auth == true {
		w.Header().Set("Authorization", fmt.Sprintf("%s %s", b.Key, security.CalcAuthHash(b.Sign, b.Date, b.URL, b.Body)))
	}
	w.Header().Set("Content-Type", b.Content)
	w.Header().Set("Auth-Date", b.Date)
	w.WriteHeader(b.Status)
	w.Write(b.Body)
	c.Clients.Prometheus.IncReqCount(b.Status, r.Method, r.URL.Path)
	c.Clients.Prometheus.ObserveReqDuration(&start, r.URL.Path)

	c.Clients.Log.Sugar().Debugw("request processed", "status", b.Status, "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
}

// EOF
