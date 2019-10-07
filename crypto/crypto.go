package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Generating a SHA256 HMAC Hash
func Hash(key, secret string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(key))

	// Proxy result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}
