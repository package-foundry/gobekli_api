package apikey

import (
	"errors"
	"strings"
)

var (
	ErrInvalidKey       = errors.New("invalid API key")
	ErrInvalidPrefix    = errors.New("invalid prefix")
	ErrKeyTooShort      = errors.New("key is too short")
	ErrMissingSeparator = errors.New("missing prefix separator")
)

func ValidateKey(key string, expectedPrefix string, expectedSize int) error {
	if len(key) == 0 {
		return ErrInvalidKey
	}

	if expectedSize > 0 && len(key) != expectedSize {
		return errors.New("key size does not match expected size")
	}

	if !strings.HasPrefix(key, expectedPrefix+Separator) {
		return ErrInvalidPrefix
	}

	if !strings.Contains(key, Separator) {
		return ErrMissingSeparator
	}

	parts := strings.SplitN(key, Separator, 2)
	if len(parts) < 2 {
		return ErrInvalidKey
	}

	prefix := parts[0]
	if prefix != expectedPrefix {
		return ErrInvalidPrefix
	}

	payload := parts[1]
	if len(payload) == 0 {
		return ErrInvalidKey
	}

	return nil
}

func ExtractPrefix(key string) (string, error) {
	if !strings.Contains(key, Separator) {
		return "", ErrInvalidKey
	}

	separatorIndex := strings.Index(key, Separator)
	return key[:separatorIndex], nil
}

func ExtractPayload(key string) (string, error) {
	if !strings.Contains(key, Separator) {
		return "", ErrInvalidKey
	}

	separatorIndex := strings.Index(key, Separator)
	if separatorIndex+1 >= len(key) {
		return "", ErrInvalidKey
	}

	return key[separatorIndex+1:], nil
}

func HasValidPrefix(key string, validPrefixes []string) bool {
	for _, p := range validPrefixes {
		if strings.HasPrefix(key, p+Separator) {
			return true
		}
	}
	return false
}
