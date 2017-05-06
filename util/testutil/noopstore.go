package testutil

import (
	"io"

	"github.com/leeola/fixity"
)

// NoopStore implements a fixity.Store as noops.
type NoopStore struct {
}

func (*NoopStore) Exists(string) (bool, error) {
	return false, nil
}

func (*NoopStore) Read(string) (io.ReadCloser, error) {
	return nil, nil
}

func (*NoopStore) Write([]byte) (string, error) {
	fakeHash, err := fixity.NewId()
	if err != nil {
		return "", err
	}
	return "fakehash-" + fakeHash, nil
}

func (*NoopStore) WriteHash(string, []byte) error {
	return nil
}

func (*NoopStore) List() (<-chan string, error) {
	return nil, nil
}
