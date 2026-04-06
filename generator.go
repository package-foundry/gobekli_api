package apikey

import (
	"crypto/rand"
	"errors"
	"strings"
)

var (
	ErrPrefixEmpty    = errors.New("prefix cannot be empty")
	ErrPrefixTooLong  = errors.New("prefix is too long for the specified key size")
	ErrPrefixHasSpace = errors.New("prefix cannot contain whitespace")
	ErrInvalidKeySize = errors.New("key size must be greater than prefix + separator")
)

const Separator = "_"

type Generator struct {
	prefix  string
	keySize int
}

func NewGenerator(prefix string, keySize int) (*Generator, error) {
	if err := validatePrefix(prefix, keySize); err != nil {
		return nil, err
	}

	return &Generator{
		prefix:  prefix,
		keySize: keySize,
	}, nil
}

func validatePrefix(prefix string, keySize int) error {
	prefix = strings.TrimSpace(prefix)

	if len(prefix) == 0 {
		return ErrPrefixEmpty
	}

	for _, r := range prefix {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			return ErrPrefixHasSpace
		}
	}

	minSize := len(prefix) + len(Separator) + 1
	if keySize < minSize {
		return ErrInvalidKeySize
	}

	if len(prefix)+len(Separator) >= keySize {
		return ErrPrefixTooLong
	}

	return nil
}

func (g *Generator) Generate() (string, error) {
	randomPartLen := g.keySize - len(g.prefix) - len(Separator)

	randomBytes := make([]byte, randomPartLen)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	encoded := base58Encode(randomBytes)

	if len(encoded) > randomPartLen {
		encoded = encoded[:randomPartLen]
	} else if len(encoded) < randomPartLen {
		return "", errors.New("generated key is shorter than expected")
	}

	return g.prefix + Separator + encoded, nil
}

func (g *Generator) Prefix() string {
	return g.prefix
}

func (g *Generator) KeySize() int {
	return g.keySize
}

func DefaultKeySize() int {
	return 32
}
