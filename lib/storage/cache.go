package storage

// Cache provides an interface to Cache storage
type Cache interface {
	Del(key []byte) (affected bool)
	Get(key []byte) (value []byte, err error)
	Set(key, value []byte, expireSeconds int) (err error)
}

// EOF
