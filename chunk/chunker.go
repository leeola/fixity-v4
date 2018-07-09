package chunk

import "context"

// Chunker implements chunking over bytes.
type Chunker interface {
	Chunk(context.Context) (Chunk, error)
}

type Chunk struct {
	Bytes []byte
	Size  int64
}
