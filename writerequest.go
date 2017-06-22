package fixity

import (
	"io"
	"os"
)

const (

	// DefaultMinChunkSize is the minimum size for rolling a blob
	//
	// See WriteRequest.MinChunkSize documentation for further explanation of
	// chunk sizes.
	//
	// This value is 1024*12 bytes, 12KiB, chosen as a reasonable size for the
	// lower allowed chunk size of common modifications.
	// Ideally it should strike a balance between too many chunks written while
	// still mitigating data duplication caused by frequent changes to files
	// and written data.
	//
	// If written data is commonly editing the first X bytes, this value should
	// be overridden in the WriteRequest itself.
	//
	// IMPORTANT: Implementations of chunkers may not make use of both min and
	// max chunk sizes. Read the implementations specific documentation to ensure
	// chunking values will work as expected.
	DefaultMinChunkSize int64 = 12288

	// DefaultMaxChunkSize is the maximum size for rolling a blob.
	//
	// See WriteRequest.MaxChunkSize documentation for further explanation of
	// chunk sizes.
	//
	// This value is 1024*12 bytes, 12KiB, chosen as a reasonable size for the
	// upper allowed chunk size of infrequent modifications.
	// Ideally it should strike a balance between too many chunks written while
	// still mitigating data duplication caused by frequent changes to files
	// and written data.
	//
	// If written data is commonly editing the first X bytes, this value should
	// be overridden in the WriteRequest itself.
	//
	// IMPORTANT: Implementations of chunkers may not make use of both min and
	// max chunk sizes. Read the implementations specific documentation to ensure
	// chunking values will work as expected.
	DefaultMaxChunkSize int64 = 4194304

	// autoChunkCount is the number of chunks SetRollFromBytes/etc will set to.
	//
	// Eg, if SetRollFromBytes is given a byte array that is 1,048,576 bytes,
	// then it will divide 1,048,576 by eg, 10 to set the rollSize to
	// 104,857, effectively rolling the bytes into 10 parts when it's written.
	autoChunkCount = 6
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
		RollSize: DefaultMaxChunkSize,
		Fields:   f,
		Blob:     rc,
	}
}

func (req *WriteRequest) setRollSizeWithMin(roll int64) {
	if roll < DefaultMinChunkSize {
		roll = DefaultMinChunkSize
	}
	if roll > DefaultMaxChunkSize {
		roll = DefaultMaxChunkSize
	}
	req.RollSize = roll
}

func (req *WriteRequest) SetRollFromBytes(b []byte) {
	req.setRollSizeWithMin(int64(len(b)) / autoChunkCount)
}

func (req *WriteRequest) SetRollFromFileInfo(fi os.FileInfo) {
	req.setRollSizeWithMin(fi.Size() / autoChunkCount)
}
