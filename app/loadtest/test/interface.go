package test

import "github.com/zk101/nixie/app/loadtest/data"

// Profile is an interface that produces a basis for creating test profiles.
type Profile interface {
	GetConfig() *Config
	GetTests() *[]Data
	Setup(*Config) error
	Start(*data.Local) error
	Stop(*data.Local) error
	Reset(*data.Local) error
}

// EOF
