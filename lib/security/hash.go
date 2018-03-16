package security

import (
	"encoding/base64"
	"strconv"

	"golang.org/x/crypto/sha3"
)

// CalcAuthHash calculates a hash and returns it
func CalcAuthHash(secret, date, url string, body []byte) string {
	hash := make([]byte, 136)
	hashStr := secret + date + url + string(body)
	sha3.ShakeSum256(hash, []byte(hashStr))
	return base64.StdEncoding.EncodeToString(hash)
}

// CalcPBhash calculates a hash and returns it
func CalcPBhash(secret, key, track string, msgType, msgSec int, date int64, body []byte) string {
	hash := make([]byte, 136)
	hashStr := secret + key + track + strconv.Itoa(msgType) + strconv.Itoa(msgSec) + strconv.FormatInt(date, 10) + string(body)
	sha3.ShakeSum256(hash, []byte(hashStr))
	return base64.StdEncoding.EncodeToString(hash)
}

// EOF
