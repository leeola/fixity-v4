package crypto

import (
	"errors"

	"github.com/leeola/kala/store/crypto/aes"
)

// Cryptoer encrypts and decrypts data depending on the backend implementation.
//
// It provides a standard interface for the stores to ignore the crypto
// implementation.
type Cryptoer interface {
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
}

type Config struct {
	CreateMissingKey bool
	KeyPath          string

	UseAes bool
}

func (c Config) UsesCrypto() bool {
	switch {
	case c.UseAes:
	default:
		return false
	}
	return true
}

func New(c Config) (Cryptoer, error) {
	switch {
	case c.UseAes:
		return aes.New(aes.Config{KeyPath: c.KeyPath})
	default:
		return nil, errors.New("No crypto specified")
	}
}
