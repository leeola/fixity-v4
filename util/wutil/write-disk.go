package wutil

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/dchest/blake2b"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/blobstore"
	"github.com/leeola/fixity/chunk"
)

const partSize = 2 // low for testing

func WriteData(ctx context.Context, w blobstore.Writer,
	chunkRefs []fixity.Ref, totalSize int64, contentHash string) ([]fixity.Ref, error) {

	chunkRefLen := len(chunkRefs)

	// -1 ensures that the morePartCount doesn't increase at an equal divide,
	// like 2 items for a pagesize of 2, would only need 1 page, yet morePartCount
	// would indicate that there's a morePart page as well.
	//
	//
	// I feel like that made no sense.
	morePartCount := (chunkRefLen - 1) / partSize

	var lastPart *fixity.Ref

	// write all of the parts first, including the partial final part..
	// ie, the part that has less than the max chunks.
	for i := morePartCount; i > 0; i-- {
		startBound := partSize * i
		endBound := startBound + partSize
		if i == morePartCount {
			endBound = startBound + chunkRefLen%partSize
		}

		part := fixity.Parts{
			Schema: fixity.Schema{
				SchemaType: fixity.BlobTypeParts,
			},
			Parts:     chunkRefs[startBound:endBound],
			MoreParts: lastPart,
		}

		ref, err := MarshalAndWrite(ctx, w, part)
		if err != nil {
			return nil, fmt.Errorf("marshalandwrite part %d: %v", i, err)
		}
		chunkRefs = append(chunkRefs, ref)
	}

	endBound := partSize
	if chunkRefLen < partSize {
		endBound = chunkRefLen
	}

	// now we've written all the parts except for the most important
	// one, the content which has a part embedded.
	data := fixity.Data{
		Parts: fixity.Parts{
			Schema: fixity.Schema{
				SchemaType: fixity.BlobTypeData,
			},
			Parts:     chunkRefs[0:endBound],
			MoreParts: lastPart,
		},
		Size: totalSize,
	}

	ref, err := MarshalAndWrite(ctx, w, data)
	if err != nil {
		return nil, fmt.Errorf("marshalandwrite content: %v", err)
	}

	return append(chunkRefs, ref), nil
}

func WriteChunks(ctx context.Context, w blobstore.Writer, r chunk.Chunker) (
	refs []fixity.Ref, totalSize int64, contentHash string, err error) {

	hasher := blake2b.New256()

	var hashes []fixity.Ref
	for {
		c, err := r.Chunk(ctx)
		if err != nil && err != io.EOF {
			return nil, 0, "", fmt.Errorf("chunk: %v", err)
		}

		totalSize += c.Size

		if err == io.EOF {
			break
		}

		if _, err := hasher.Write(c.Bytes); err != nil {
			return nil, 0, "", fmt.Errorf("hasher write: %v", err)
		}

		h, err := w.Write(ctx, c.Bytes)
		if err != nil {
			return nil, 0, "", fmt.Errorf("blob write: %v", err)
		}

		hashes = append(hashes, h)
	}

	hash := hex.EncodeToString(hasher.Sum(nil)[:])
	return hashes, totalSize, hash, nil
}

func MarshalAndWrite(ctx context.Context, w blobstore.Writer, v interface{}) (fixity.Ref, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("marshal: %v", err)
	}

	ref, err := w.Write(ctx, b)
	if err != nil {
		return "", fmt.Errorf("blob write: %v", err)
	}

	return ref, nil
}

func WriteValues(ctx context.Context, w blobstore.Writer, v fixity.Values) (fixity.Ref, error) {
	return MarshalAndWrite(ctx, w, fixity.ValuesMap{
		Schema: fixity.Schema{
			SchemaType: fixity.BlobTypeValues,
		},
		Values: v,
	})
}
