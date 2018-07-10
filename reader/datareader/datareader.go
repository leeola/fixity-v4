package datareader

import (
	"context"
	"fmt"
	"io"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore"
)

type Reader struct {
	ctx     context.Context
	bs      blobstore.Reader
	dataRef fixity.Ref

	partReadCloser          io.ReadCloser
	parts                   []fixity.Ref
	partsIndex, partsLength int
	nextPartsRef            *fixity.Ref

	data fixity.DataSchema
}

func New(ctx context.Context, bs blobstore.Reader, ref fixity.Ref) (*Reader, error) {
	return &Reader{
		ctx:     ctx,
		bs:      bs,
		dataRef: ref,
	}, nil
}

func (r *Reader) dataStruct() error {
	var data fixity.DataSchema
	if err := blobstore.ReadAndUnmarshal(r.ctx, r.bs, r.dataRef, &data); err != nil {
		return fmt.Errorf("readandunmarshal %q: %v", r.dataRef, err)
	}

	partsLength := len(data.PartsSchema.Parts)
	if partsLength == 0 {
		return fmt.Errorf("dataschema %q missing parts", r.dataRef)
	}

	r.parts = data.PartsSchema.Parts
	r.nextPartsRef = data.MoreParts

	firstPart := data.PartsSchema.Parts[0]
	rc, err := r.bs.Read(r.ctx, firstPart)
	if err != nil {
		return fmt.Errorf("dataschema %q read: %v", r.dataRef, err)
	}

	r.partReadCloser = rc
	r.partsIndex++
	r.partsLength = partsLength
	r.data = data

	return nil
}

func (r *Reader) nextParts() error {
	if r.nextPartsRef == nil {
		return io.EOF
	}

	var parts fixity.PartsSchema
	if err := blobstore.ReadAndUnmarshal(r.ctx, r.bs, *r.nextPartsRef, &parts); err != nil {
		return fmt.Errorf("readandunmarshal: %v", err)
	}

	partsLength := len(parts.Parts)
	if partsLength == 0 {
		return fmt.Errorf("partschema %q missing parts", r.nextPartsRef)
	}

	r.partsIndex = 0
	r.partsLength = partsLength
	r.parts = parts.Parts
	r.nextPartsRef = parts.MoreParts

	return nil
}

func (r *Reader) nextPart() error {
	// close the previous part if we're trying to load
	// the next part.
	if err := r.partReadCloser.Close(); err != nil {
		return fmt.Errorf("close part: %v", err)
	}

	if r.partsIndex == r.partsLength {
		err := r.nextParts()
		if err == io.EOF {
			return io.EOF
		}
		if err != nil {
			return fmt.Errorf("nextparts: %v", err)
		}
	}

	ref := r.parts[r.partsIndex]
	rc, err := r.bs.Read(r.ctx, ref)
	if err != nil {
		return fmt.Errorf("read %q: %v", ref, err)
	}

	r.partReadCloser = rc
	r.partsIndex++

	return nil
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.partReadCloser == nil {
		if err := r.dataStruct(); err != nil {
			return 0, fmt.Errorf("dataschema: %v", err)
		}
	}

	n, err := r.partReadCloser.Read(p)
	if err == io.EOF {
		err := r.nextPart()
		if err == io.EOF {
			return n, io.EOF
		}
		if err != nil {
			return 0, fmt.Errorf("nextpart: %v", err)
		}
		return n, nil
	}
	return n, err
}

func (r *Reader) Checksum() (string, error) {
	if r.partReadCloser == nil {
		if err := r.dataStruct(); err != nil {
			return "", fmt.Errorf("dataschema: %v", err)
		}
	}

	return r.data.Checksum, nil
}

func (r *Reader) Size() (int64, error) {
	if r.partReadCloser == nil {
		if err := r.dataStruct(); err != nil {
			return 0, fmt.Errorf("dataschema: %v", err)
		}
	}

	return r.data.Size, nil
}
