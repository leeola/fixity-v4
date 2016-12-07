package store

import (
	"encoding/json"
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

func WriteContentRoller(s Store, r ContentRoller) ([]string, error) {
	var hashes []string
	for {
		c, err := r.Roll()
		if err != nil && err != io.EOF {
			return nil, errors.Stack(err)
		}

		if err == io.EOF {
			break
		}

		h, err := WriteContent(s, c)
		if err != nil {
			return nil, errors.Stack(err)
		}
		hashes = append(hashes, h)
	}

	return hashes, nil
}

func WriteMultiPart(s Store, mp MultiPart) (string, error) {
	b, err := json.Marshal(mp)
	if err != nil {
		return "", errors.Stack(err)
	}

	h, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func WriteContent(s Store, c Content) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", errors.Stack(err)
	}

	h, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func MultiPartFromReader(io.Reader) (MultiPart, error) {
	return MultiPart{}, errors.New("not implemented")
}
