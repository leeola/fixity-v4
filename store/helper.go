package store

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
	"github.com/leeola/kala/util"
)

// ContentType exists to get the content type from a blob.
//
// Note that this only matches the contentType field, and only applies to the Meta
// struct (or anything defining contentType).
type ContentType struct {
	ContentType string `json:"contentType"`
}

func NewAnchor(s Store) (string, error) {
	return WriteAnchor(s, Anchor{
		AnchorRand: util.RandomInt(),
	})
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

func WriteAnchor(s Store, a Anchor) (string, error) {
	return MarshalAndWrite(s, a)
}

func WriteMultiPart(s Store, mp MultiPart) (string, error) {
	return MarshalAndWrite(s, mp)
}

func WriteMeta(s Store, m Meta) (string, error) {
	return MarshalAndWrite(s, m)
}

func WriteContent(s Store, c Content) (string, error) {
	return MarshalAndWrite(s, c)
}

func MultiPartFromReader(io.Reader) (MultiPart, error) {
	return MultiPart{}, errors.New("not implemented")
}

func ReadAndUnmarshal(s Store, h string, v interface{}) error {
	rc, err := s.Read(h)
	if err != nil {
		return errors.Stack(err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return errors.Stack(err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.Stack(err)
	}

	return nil
}

func ReadMeta(s Store, h string) (Meta, error) {
	var m Meta
	if err := ReadAndUnmarshal(s, h, &m); err != nil {
		return Meta{}, errors.Stack(err)
	}

	if !IsValidMeta(m) {
		return Meta{}, errors.Errorf("given hash %q is not a valid meta struct", h)
	}

	return m, nil
}

func GetContentTypeWithReader(s Store, h string) (string, io.ReadCloser, error) {
	rc, err := s.Read(h)
	if err != nil {
		return "", nil, errors.Stack(err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", nil, errors.Stack(err)
	}

	var ct ContentType
	if err := json.Unmarshal(b, &ct); err != nil {
		return "", nil, errors.Stack(err)
	}

	return ct.ContentType, ioutil.NopCloser(bytes.NewReader(b)), nil
}

func IsValidMeta(m Meta) bool {
	return m.UploadedAt.IsZero() && (m.Anchor == "" || m.Multi == "")
}
