package util

import (
	"encoding/base64"
	"fmt"
)

// B64Decode decodes the given base64-encoded string.
func B64Decode(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		panic(fmt.Sprintf("Could not decode base64: %v", err))
	}
	return string(decoded)
}
