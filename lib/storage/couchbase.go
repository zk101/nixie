package storage

import (
	"github.com/couchbase/gocb"
)

// Couchbase provides an interface for Couchbase Server
type Couchbase interface {
	Get(key string, valuePtr interface{}) (gocb.Cas, error)
	Insert(key string, value interface{}, expiry uint32) (gocb.Cas, error)
	Remove(key string, cas gocb.Cas) (gocb.Cas, error)
	Replace(key string, value interface{}, cas gocb.Cas, expiry uint32) (gocb.Cas, error)
	Touch(key string, cas gocb.Cas, expiry uint32) (gocb.Cas, error)
	Upsert(key string, value interface{}, expiry uint32) (gocb.Cas, error)
}

// EOF
