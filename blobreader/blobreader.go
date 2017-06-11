package blobreader

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/leeola/errors"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/store"
)

// HashWithBytes is used by Reader to return the data of a hash if needed.
type HashWithBytes struct {
	Hash  string
	Bytes []byte
}

type Config struct {
	Hash  string
	Store store.Store
}

type Reader struct {
	hashes []string
	store  store.Store

	currentReader io.Reader
}

func New(s fixity.Store, blobHash string) (*Reader, error) {
	return &Reader{
		hashes: []string{blobHash},
		store:  s,
	}, nil
}

func (r *Reader) Read(p []byte) (int, error) {
	n, hwb, err := r.ReadContentOnly(p)
	if err == io.EOF {
		return 0, io.EOF
	}
	if err != nil {
		return 0, errors.Stack(err)
	}

	// okay if nil, will noop
	if err := r.UnmarshalHashes(hwb.Bytes); err != nil {
		return 0, errors.Stack(err)
	}

	return n, nil
}

// ReadContentOnly returns the read bytes if they do not contain store.Content.
//
// This allows another Reader (eg: indexreader.Reader) to create a reader for
// additional structures, such as querying the index for an anchor and adding
// it to the store.Reader.
func (r *Reader) ReadContentOnly(p []byte) (int, HashWithBytes, error) {
	if r.currentReader != nil {
		n, err := r.currentReader.Read(p)
		if err == io.EOF {
			r.currentReader = nil
		} else if err != nil {
			return 0, HashWithBytes{}, errors.Stack(err)
		}

		return n, HashWithBytes{}, nil
	}

	if len(r.hashes) <= 0 {
		return 0, HashWithBytes{}, io.EOF
	}

	// pop the first hash, as that has read priority.
	h := r.hashes[0]
	r.hashes = r.hashes[1:]

	// Load the hash and unmarshal it.
	rc, err := r.store.Read(h)
	if err != nil {
		return 0, HashWithBytes{}, errors.Stack(err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return 0, HashWithBytes{}, errors.Stack(err)
	}

	return 0, HashWithBytes{
		Hash:  h,
		Bytes: b,
	}, nil
}

func (r *Reader) UnmarshalHashes(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	var d readerData
	if err := json.Unmarshal(b, &d); err != nil {
		return errors.Stack(err)
	}

	switch {
	case d.Meta != "":
		r.AddHashes(d.Meta)

	case d.MultiPart != "":
		r.AddHashes(d.MultiPart)

	case len(d.Parts) > 0:
		r.AddHashes(d.Parts...)

	case len(d.Part) > 0:
		r.SetCurrentReader(bytes.NewReader(d.Part))

	default:
		return errors.New("Reader: unhandled hash content")
	}

	return nil
}

func (r *Reader) AddHashes(s ...string) {
	r.hashes = append(r.hashes, s...)
}

func (r *Reader) SetCurrentReader(cr io.Reader) {
	r.currentReader = cr
}
