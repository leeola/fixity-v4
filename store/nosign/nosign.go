package nosign

import (
	"context"
	"errors"
	"io"
)

type Store struct {
	path string
}

func New(path string) (*Store, error) {
	return &Store{path: path}, nil
}

func (s *Store) NewID(_ context.Context) (string, error) {
	return "", errors.New("not implemented")
}

func (s *Store) Write(_ context.Context, id string, r io.Reader) ([]string, error) {
	return nil, errors.New("not implemented")
}
