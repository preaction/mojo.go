package util

import (
	"crypto/md5"
	"encoding/base32"
)

// MD5Sum returns a base32-encoded string of the MD5 sum of the input
// string.
func MD5Sum(input string) string {
	sum := md5.Sum([]byte(input))
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sum[:])
}
