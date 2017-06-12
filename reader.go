package fixity

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/fatih/structs"
	"github.com/leeola/errors"
)

type BlobReader struct {
	blobHash      string
	loadedChunks  bool
	hashes        []string
	store         Store
	currentReader io.ReadCloser
}

func Reader(s Store, blobHash string) *BlobReader {
	return &BlobReader{
		blobHash: blobHash,
		store:    s,
	}
}

func (r *BlobReader) Read(p []byte) (int, error) {
	if !r.loadedChunks {
		if err := r.loadChunks(); err != nil {
			return 0, err
		}
	}

	if len(r.hashes) == 0 {
		return 0, io.EOF
	}

	if r.currentReader == nil {
		h := r.hashes[0]
		r.hashes = r.hashes[1:]
		var c Chunk
		if err := readAndUnmarshal(r.store, h, &c); err != nil {
			return 0, err
		}
		r.currentReader = ioutil.NopCloser(bytes.NewReader(c.ChunkBytes))
	}

	i, err := r.currentReader.Read(p)
	if err == io.EOF {
		r.currentReader = nil
		err = nil
	}
	return i, err
}

func (r *BlobReader) loadChunks() error {
	var b Blob
	if err := readAndUnmarshal(r.store, r.blobHash, &b); err != nil {
		return err
	}

	if structs.IsZero(b) {
		return errors.New("blobHash is not valid Blob type")
	}

	r.loadedChunks = true
	r.hashes = b.ChunkHashes

	return nil
}

func (r *BlobReader) Close() error {
	r.hashes = nil
	if r.currentReader == nil {
		return nil
	}
	return r.currentReader.Close()
}
