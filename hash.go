package apikey

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashKey(key string) []byte {
	h := sha256.Sum256([]byte(key))
	return h[:]
}

func HashKeyHex(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func Fingerprint(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:8])
}

func FingerprintBytes(key string) []byte {
	h := sha256.Sum256([]byte(key))
	return h[:8]
}
