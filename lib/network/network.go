package network

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

// CheckHTTPhost tests a r.Remote against a pattern and returns an error
func CheckHTTPhost(host, pattern string) error {
	split := strings.Split(host, ":")
	if len(split) != 2 {
		return errors.New("split on host string failed")
	}
	matched, err := regexp.MatchString(pattern, host)
	if err != nil {
		return err
	}
	if matched != true {
		return errors.New("host match failed")
	}

	return nil
}

// FindInetAddress finds and returns the first match interface address
func FindInetAddress(pattern string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			matched, err := regexp.MatchString(pattern, addr.String())
			if err != nil {
				return "", err
			}
			if matched == true {
				split := strings.Split(addr.String(), "/")
				if len(split) != 2 {
					return "", errors.New("split failed")
				}
				return split[0], nil
			}
		}
	}

	return "", errors.New("no match found")
}

// EOF
