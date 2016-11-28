package store

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
)

func NewPerma(s Store, h string) (string, error) {
	return "", errors.New("not implemented")
}

func WriteReader(s Store, r io.Reader) (string, error) {
	if s == nil {
		return "", errors.New("Store is nil")
	}
	if r == nil {
		return "", errors.New("Reader is nil")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to readall")
	}

	fmt.Printf("bytes: %q\n", string(b))

	h, err := s.Write(b)
	return h, errors.Wrap(err, "store failed to write")
}

func WriteHashReader(s Store, h string, r io.Reader) error {
	if s == nil {
		return errors.New("Store is nil")
	}
	if r == nil {
		return errors.New("Reader is nil")
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return errors.Wrap(err, "failed to readall")
	}

	err = s.WriteHash(h, b)
	// do not wrap hash not match err
	if err == HashNotMatchContentErr {
		return err
	}
	return errors.Wrap(err, "store failed to write")
}

func WriteMultiPart(s Store, c MultiPart) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func WriteContent(s Store, c Content) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}

func ReadParts(s Store, h string) (io.ReadCloser, error) {
	return nil, errors.New("not implemented")
}
