package store

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
)

type readerData struct {
	Parts   []string `json:"parts"`
	Content []byte   `json:"content"`
}

type ReaderConfig struct {
	Hash  string
	Store Store
}

type Reader struct {
	hashes []string
	store  Store

	currentReader io.Reader
}

func NewReader(c ReaderConfig) (*Reader, error) {
	if c.Hash == "" {
		return nil, errors.New("missing required config field: Hash")
	}
	if c.Store == nil {
		return nil, errors.New("missing required config field: Store")
	}

	return &Reader{
		hashes: []string{c.Hash},
		store:  c.Store,
	}, nil
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.currentReader != nil {
		n, err := r.currentReader.Read(p)
		if err == io.EOF {
			r.currentReader = nil
		} else if err != nil {
			return 0, errors.Stack(err)
		}

		return n, nil
	}

	if len(r.hashes) <= 0 {
		return 0, io.EOF
	}

	// pop the first hash, as that has read priority.
	h := r.hashes[0]
	r.hashes = r.hashes[1:]

	// Load the hash and unmarshal it.
	rc, err := r.store.Read(h)
	if err != nil {
		return 0, errors.Stack(err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return 0, errors.Stack(err)
	}

	var d readerData
	if err := json.Unmarshal(b, &d); err != nil {
		return 0, errors.Stack(err)
	}

	if len(d.Parts) > 0 {
		r.hashes = append(r.hashes, d.Parts...)
	}

	if len(d.Content) > 0 {
		r.currentReader = bytes.NewReader(d.Content)
	}

	return 0, nil
}
