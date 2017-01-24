package aes

import (
	"encoding/base64"
	"io/ioutil"
	"os"
)

type Config struct {
	KeyPath string
}

type Aes struct {
	Key *[32]byte
}

func New(c Config) (*Aes, error) {
	key, err := LoadOrNew(c.KeyPath)
	if err != nil {
		return nil, err
	}

	return &Aes{
		Key: key,
	}, nil
}

func (a *Aes) Decrypt(b []byte) ([]byte, error) {
	return Decrypt(b, a.Key)
}

func (a *Aes) Encrypt(b []byte) ([]byte, error) {
	return Encrypt(b, a.Key)
}

func LoadOrNew(p string) (*[32]byte, error) {
	k, err := LoadKey(p)
	if os.IsNotExist(err) {
		return NewKey(p)
	}
	return k, err
}

func LoadKey(p string) (*[32]byte, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	d := make([]byte, len(b))
	n, err := base64.StdEncoding.Decode(d, b)
	if err != nil {
		return nil, err
	}

	var key [32]byte
	copy(key[:], d[:n])

	return &key, nil
}

// NewKey generates and then saved a key to the given path
func NewKey(p string) (*[32]byte, error) {
	k := NewEncryptionKey()

	s := base64.StdEncoding.EncodeToString(k[:])
	if err := ioutil.WriteFile(p, []byte(s), 0644); err != nil {
		return nil, err
	}
	return k, nil
}
