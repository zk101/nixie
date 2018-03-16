package config

import (
	"errors"
	"os"

	"github.com/zk101/nixie/lib/config"
	"github.com/zk101/nixie/lib/consul"
)

// Load processes cli arguments and figures out config
func Load() (*Config, error) {
	flags := processFlags()

	if flags.configFile == "" && flags.consulAddress == "" {
		return nil, errors.New("require a config file or consul address flag")
	}

	conf := NewConfig()

	if flags.configFile != "" {
		if err := config.Load(flags.configFile, conf); err != nil {
			return nil, err
		}
	} else if flags.consulAddress != "" {
		config.Override(conf, flags.consulPrefix, "/", consul.NewClient(&consul.Config{Address: flags.consulAddress, Scheme: flags.consulScheme, Datacenter: flags.consulDatacenter, Token: flags.consulToken}))
	}

	conf.loadOverrides(consul.NewClient(&conf.Consul))

	if conf.Controls.SSL == true {
		if conf.Controls.SSLkey == "" {
			return nil, errors.New("no sslkey specified")
		}

		if _, err := os.Stat(conf.Controls.SSLkey); err != nil {
			return nil, err
		}

		if conf.Controls.SSLcert == "" {
			return nil, errors.New("no sslcert specified")
		}

		if _, err := os.Stat(conf.Controls.SSLcert); err != nil {
			return nil, err
		}
	}

	return conf, nil
}

// loadOverrides processes Consul and Env overrides
func (c *Config) loadOverrides(consul *consul.Client) {
	if c.Controls.OverrideConsul == true {
		config.Override(c, c.Controls.PrefixConsul, "/", consul)
	}

	if c.Controls.OverrideEnv == true {
		config.Override(c, c.Controls.PrefixEnv, "_", &config.KVenv{})
	}
}

// EOF
