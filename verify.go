package apikey

import (
	"crypto/subtle"
)

func VerifyKey(presented string, storedHash []byte) bool {
	presentedHash := HashKey(presented)
	return subtle.ConstantTimeCompare(presentedHash, storedHash) == 1
}

func VerifyKeyHex(presented string, storedHashHex string) bool {
	presentedHash := HashKeyHex(presented)
	return subtle.ConstantTimeCompare(
		[]byte(presentedHash),
		[]byte(storedHashHex),
	) == 1
}
