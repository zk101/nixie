package testutil

import (
	"errors"
	"os"

	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/couchbase"
	"github.com/zk101/nixie/lib/ldap"
)

// Config holds common structs used in testing
type Config struct {
	Controls  Controls
	Couchbase couchbase.Config
	LDAP      ldap.Config
}

// Controls holds application controls used to support testing
type Controls struct {
	CAcertPath string
}

// NewConfig returns a Config struct ptr with DefaultConfig loaded
func NewConfig() *Config {
	return &Config{
		Couchbase: couchbase.DefaultConfig(),
		LDAP:      ldap.DefaultConfig(),
	}
}

// LoadConfig attempts to load the testutil config file and returns a ptr to Config
func LoadConfig() (*Config, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return nil, errors.New("GOPATH must be set")
	}

	loadFile := "dev.toml"
	envFile := os.Getenv("TESTUTIL_CONFIGFILE")
	if envFile != "" {
		loadFile = envFile
	}
	configFile := gopath + "/src/github.com/zk101/nixie/etc/testutil/" + loadFile

	conf := NewConfig()
	if err := config.Load(configFile, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// EOF
