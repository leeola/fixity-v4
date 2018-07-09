package blobreader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/leeola/fixity"
)

// ReadCloser implements peek/buffer into a BlobType to identify the blob,
// and then caches the result for repeated BlobType requests.
type ReadCloser struct {
	*bytes.Buffer
	blobType fixity.BlobType
}

func BlobType(r io.Reader) (*ReadCloser, fixity.BlobType, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, 0, fmt.Errorf("readall: %v", err)
	}

	// TODO(leeola): determine if unmarshalling large blobs of non-json
	// is inefficient compared to peeking the first byte. If first byte
	// peeking is better, add it above this section.

	// if the data fails to unmarshal, consider it schemaless.
	// the zero value of fixity.Schema.SchemaType == fixity.BlobTypeSchemaless
	var schema fixity.Schema
	_ = json.Unmarshal(b, &schema)

	return &ReadCloser{
		Buffer:   bytes.NewBuffer(b),
		blobType: schema.SchemaType,
	}, schema.SchemaType, nil
}

func (rc *ReadCloser) BlobType() fixity.BlobType {
	return rc.blobType
}

func (rc *ReadCloser) Close() error {
	return nil
}
