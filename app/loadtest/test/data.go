package test

import "github.com/zk101/nixie/app/loadtest/data"

// Data provides the struct to build individual tests
type Data struct {
	ID       string
	Callback func(*data.Local) error
	Count    uint
}

// EOF
