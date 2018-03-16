package security

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
)

// ldapSalt returns a random string suitable for LDAP
func ldapSalt(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return []byte(""), err
	}
	return bytes, nil
}

// LDAPhashSSHA512 returns a string compatable to insert into an ldap password attribute
func LDAPhashSSHA512(passwd string) (string, error) {
	salt, err := ldapSalt(8)
	if err != nil {
		return "", err
	}

	data := sha512.Sum512(append([]byte(passwd), salt[:]...))

	return "{ssha512}" + base64.StdEncoding.EncodeToString(append(data[:], salt[:]...)), nil
}

// EOF
