package config

import (
	"os"
	"strings"
)

// KVenv provides struct
type KVenv struct{}

// Get checks env and returns a value
func (k *KVenv) Get(key string) (string, error) {
	value := os.Getenv(strings.ToUpper(key))

	return value, nil
}

// EOF
