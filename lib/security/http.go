package security

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// ValidateHeaders checks headers are appropriate for request
func ValidateHeaders(r *http.Request) error {
	if r.URL.RawQuery != "" {
		return fmt.Errorf("found query string %s", r.URL.RawQuery)
	}

	if r.Proto != "HTTP/1.0" && r.Proto != "HTTP/1.1" && r.Proto != "HTTP/2" {
		return fmt.Errorf("found http protocol version %s", r.Proto)
	}

	if len(r.Header["Authorization"]) > 1 {
		return fmt.Errorf("found %d authorization headers", len(r.Header["Authorization"]))
	}

	if len(r.Header["Auth-Date"]) > 1 {
		return fmt.Errorf("found %d auth date headers", len(r.Header["Authorization"]))
	}

	if len(r.Header["Accept"]) != 1 {
		return fmt.Errorf("found %d accept headers", len(r.Header["Accept"]))
	}

	if len(r.Header["Accept-Encoding"]) != 1 {
		return fmt.Errorf("found %d accept encoding headers", len(r.Header["Accept-Encoding"]))
	}

	if r.Header["Accept-Encoding"][0] != "gzip" {
		return errors.New("accept encoding header does not equal gzip")
	}

	if len(r.Header["User-Agent"]) != 1 {
		return fmt.Errorf("found %d user agent headers", len(r.Header["User-Agent"]))
	}

	switch r.Header["User-Agent"][0] {
	case "Consul Health Check":
		if r.Header["Accept"][0] != "text/plain, text/*, */*" {
			return fmt.Errorf("found unexpected accept header %s", r.Header["Accept"][0])
		}

	case "loadtest":
		if r.Header["Accept"][0] != "application/octet-stream" {
			return fmt.Errorf("found unexpected accept header %s", r.Header["Accept"][0])
		}

		if len(r.Header["Content-Type"]) != 1 {
			return fmt.Errorf("found %d content type headers", len(r.Header["Content-Type"]))
		}

		if r.Header["Content-Type"][0] != "application/octet-stream" {
			return fmt.Errorf("found unexpected accept header %s", r.Header["Content-Type"][0])
		}

	default:
		matched, err := regexp.MatchString(`^Prometheus\/\d\.\d\.\d$`, r.Header["User-Agent"][0])
		if err != nil {
			return fmt.Errorf("found unexpected accept header %s", err.Error())
		}

		if matched == true {
			switch r.Header["Accept"][0] {
			case "application/vnd.google.protobuf;proto=io.prometheus.client.MetricFamily;encoding=delimited;q=0.7,text/plain;version=0.0.4;q=0.3,*/*;q=0.1":

			case "text/plain;version=0.0.4;q=1,*/*;q=0.1":

			default:
				return fmt.Errorf("found unexpected accept header %s", r.Header["Accept"][0])
			}

			if len(r.Header["X-Prometheus-Scrape-Timeout-Seconds"]) != 1 {
				return fmt.Errorf("found %d x prometheus scrape timeout seconds headers", len(r.Header["X-Prometheus-Scrape-Timeout-Seconds"]))
			}
		} else {
			return fmt.Errorf("found unexpected user agent %s", r.Header["User-Agent"][0])
		}

		return nil
	}

	if len(r.Header["Connection"]) != 1 {
		return fmt.Errorf("found %d connection headers", len(r.Header["Connection"]))
	}

	if r.Header["Connection"][0] != "close" {
		return errors.New("connection header does not equal close")
	}

	return nil
}

// ValidateAuthorization validates the authorization header from a request
func ValidateAuthorization(r *http.Request, t uint32) (string, string, string, error) {
	if len(r.Header["Authorization"]) != 1 {
		return "", "", "", errors.New("Require only one Authorization Header")
	}

	if len(r.Header["Auth-Date"]) != 1 {
		return "", "", "", errors.New("Require only one Auth-Date Header")
	}

	authBits := strings.Split(r.Header["Authorization"][0], " ")
	if len(authBits) != 2 {
		return "", "", "", errors.New("Authorization Header does not consist of 2 parts")
	}

	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{32}$`, authBits[0]); matched == false {
		return "", "", "", errors.New("Access Key contains invalid characters")
	}

	if matched, _ := regexp.MatchString(`^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`, authBits[1]); matched == false {
		return "", "", "", errors.New("Access Hash not a base64 encoded string")
	}

	if matched, _ := regexp.MatchString(`^(Mon|Tue|Wed|Thu|Fri|Sat|Sun), \d{1,2} (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{4} \d{2}:\d{2}:\d{2} [A-Z]{3}$`, r.Header["Auth-Date"][0]); matched == false {
		return "", "", "", errors.New("Invalid characters found in Date")
	}

	timeReceived, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", r.Header["Auth-Date"][0])
	if err != nil {
		return "", "", "", errors.New("Unix_t did not convert to 64bit number")
	}

	if timeReceived.Before(time.Now().Add(-time.Duration(t)*time.Second)) == true || timeReceived.After(time.Now().Add(time.Duration(t)*time.Second)) == true {
		return "", "", "", errors.New("Date outside acceptable range")
	}

	return authBits[0], authBits[1], r.Header["Auth-Date"][0], nil
}

// EOF
