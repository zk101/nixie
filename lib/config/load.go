package config

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

// Load reads a config file into a Config Struct
func Load(file string, config interface{}) error {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	suffix := strings.Split(file, ".")

	switch suffix[len(suffix)-1] {
	case "json":
		if err := json.Unmarshal(fileData, config); err != nil {
			return err
		}

	case "toml":
		if _, err := toml.Decode(string(fileData), config); err != nil {
			return err
		}

	case "yaml":
		if err := yaml.Unmarshal(fileData, config); err != nil {
			return err
		}

	default:
		return errors.New("unsupported config file type")
	}

	return nil
}

// LoadCAcerts adds certs from a folder to a TLSConfig struct and returns a pointer
func LoadCAcerts(caCertDir string) (*x509.CertPool, error) {
	if caCertDir == "" {
		return nil, nil
	}

	files, err := ioutil.ReadDir(caCertDir)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	foundCerts := false

	for _, file := range files {
		if file.IsDir() == false {
			caCert, err := ioutil.ReadFile(caCertDir + "/" + file.Name())
			if err == nil {
				if caCertPool.AppendCertsFromPEM(caCert) == true {
					foundCerts = true
				}
			}
		}
	}

	if foundCerts == true {
		return caCertPool, nil
	}

	return nil, nil
}

// EOF
