package httpd

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zk101/nixie/lib/httpd"
)

// StartHTTPD sets up and starts the httpd server
func StartHTTPD(config *httpd.Config) (*http.Server, error) {
	http.HandleFunc("/", http.NotFound)
	http.HandleFunc("/health", handlerHealth)
	http.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr: ":" + strconv.Itoa(config.Port),
	}
	s.SetKeepAlivesEnabled(config.Keepalive)

	if config.SSL == true {
		if config.SSLkey == "" {
			return nil, errors.New("Missing SSL Key")
		}

		if _, err := os.Stat(config.SSLkey); err != nil {
			return nil, err
		}

		if config.SSLcert == "" {
			return nil, errors.New("Missing SSL Cert")
		}

		if _, err := os.Stat(config.SSLcert); err != nil {
			return nil, err
		}
	}

	go func() {
		if config.SSL == true {
			s.ListenAndServeTLS(config.SSLcert, config.SSLkey)
		} else {
			s.ListenAndServe()
		}
	}()

	return s, nil
}

// StopHTTPD sets up and starts the httpd server
func StopHTTPD(s *http.Server) error {
	return s.Close()
}

// handlerHealth deals with health requests
func handlerHealth(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "200 Okay", http.StatusOK)
}

// EOF
