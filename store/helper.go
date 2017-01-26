package store

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
)

// IsVersionWithBytes checks if the given hash is a Version, and returns the bytes.
//
// By returning the bytes, the caller can check if the content is a contentType.
func IsVersionWithBytes(s Store, h string) (bool, Version, []byte, error) {
	var v Version
	b, err := ReadAndUnmarshalWithBytes(s, h, &v)
	if err != nil {
		return false, Version{}, nil, errors.Stack(err)
	}

	return v.Meta != "", v, b, nil
}

// NewAnchor generates random bytes and returns the hex of those bytes.
func NewAnchor() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
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

func WritePartRoller(s Store, r PartRoller) ([]string, error) {
	var hashes []string
	for {
		c, err := r.Roll()
		if err != nil && err != io.EOF {
			return nil, errors.Stack(err)
		}

		if err == io.EOF {
			break
		}

		h, err := WritePart(s, c)
		if err != nil {
			return nil, errors.Stack(err)
		}
		hashes = append(hashes, h)
	}

	return hashes, nil
}

func MarshalAndWrite(s Store, v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", errors.Stack(err)
	}

	h, err := s.Write(b)
	if err != nil {
		return "", errors.Stack(err)
	}

	return h, nil
}

func WriteMeta(s Store, m Meta) (string, error) {
	return MarshalAndWrite(s, m)
}

func WriteMultiHash(s Store, mh MultiHash) (string, error) {
	return MarshalAndWrite(s, mh)
}

func WriteMultiPart(s Store, mp MultiPart) (string, error) {
	return MarshalAndWrite(s, mp)
}

func WritePart(s Store, c Part) (string, error) {
	return MarshalAndWrite(s, c)
}

func MultiPartFromReader(io.Reader) (MultiPart, error) {
	return MultiPart{}, errors.New("not implemented")
}

func ReadAll(s Store, h string) ([]byte, error) {
	rc, err := s.Read(h)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rc.Close()

	return ioutil.ReadAll(rc)
}

func ReadAndUnmarshal(s Store, h string, v interface{}) error {
	_, err := ReadAndUnmarshalWithBytes(s, h, v)
	return err
}

func ReadAndUnmarshalWithBytes(s Store, h string, v interface{}) ([]byte, error) {
	b, err := ReadAll(s, h)
	if err != nil {
		return nil, errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return nil, errors.Stack(err)
	}

	return b, nil
}

func ReadVersion(s Store, h string) (Version, error) {
	var v Version
	if err := ReadAndUnmarshal(s, h, &v); err != nil {
		return Version{}, errors.Stack(err)
	}

	return v, nil
}
