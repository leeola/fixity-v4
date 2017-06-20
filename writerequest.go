package fixity

import (
	"io"
	"os"
)

const (
	// DefaultRollSize is the minimum size for rolling a blob.
	//
	// See WriteRequest.RollSize documentation for further explanation of RollSizes.
	//
	// This value is (1024^2)*4 bytes, 4MiB, chosen as a reasonable size for
	// random data sources. Ideally it should strike a balance between too many
	// chunks written while still mitigating data duplication caused by
	// frequent changes to files.
	//
	// If a file format is commonly editing the first X bytes, this value should
	// be overridden in the WriteRequest itself.
	DefaultRollSize int64 = 4194304

	// autoChunkCount is the number of chunks SetRollFromBytes/etc will set to.
	//
	// Eg, if SetRollFromBytes is given a byte array that is 1,048,576 bytes,
	// then it will divide 1,048,576 by eg, 10 to set the rollSize to
	// 104,857, effectively rolling the bytes into 10 parts when it's written.
	autoChunkCount = 6

	// minAutoChunk is the smallest rollSize set by SetRollFromBytes/etc.
	minAutoChunk int64 = 1024

	// minAutoChunk is the largest rollSize set by SetRollFromBytes/etc.
	maxAutoChunk int64 = DefaultRollSize
)

// WriteRequest represents a blob to be written alone with metadata.
type WriteRequest struct {
	Id       string
	RollSize int64
	Fields   Fields
	Blob     io.ReadCloser
}

// NewWrite creates a new WriteRequest to be used with Fixity.WriteRequest.
func NewWrite(id string, rc io.ReadCloser, f ...Field) *WriteRequest {
	return &WriteRequest{
		Id:       id,
		RollSize: DefaultRollSize,
		Fields:   f,
		Blob:     rc,
	}
}

func (req *WriteRequest) setRollSizeWithMin(roll int64) {
	if roll < minAutoChunk {
		roll = minAutoChunk
	}
	if roll > maxAutoChunk {
		roll = DefaultRollSize
	}
	req.RollSize = roll
}

func (req *WriteRequest) SetRollFromBytes(b []byte) {
	req.setRollSizeWithMin(int64(len(b)) / autoChunkCount)
}

func (req *WriteRequest) SetRollFromFileInfo(fi os.FileInfo) {
	req.setRollSizeWithMin(fi.Size() / autoChunkCount)
}
