package fixity

import (
	"io"
	"os"
)

const (

	// DefaultAverageChunkSize that the chunker will attempt to produce.
	//
	// See WriteRequest.AverageChunkSize documentation for further explanation
	// of chunk sizes.
	//
	// This value is 1024^2*4 bytes, 4MiB, chosen as a compromise for the
	// small and large files.
	//
	// Ideally it should strike a balance between too many chunks written while
	// still mitigating data duplication caused by frequent changes to files
	// and written data.
	//
	// If written data is commonly editing the first X bytes, this value should
	// be overridden in the WriteRequest itself.
	DefaultAverageChunkSize uint64 = 4194304

	// autoChunkCount is the number of chunks SetRollFromBytes/etc will set to.
	//
	// Eg, if SetRollFromBytes is given a byte array that is 1,048,576 bytes,
	// then it will divide 1,048,576 by eg, 10 to set the rollSize to
	// 104,857, effectively rolling the bytes into 10 parts when it's written.
	autoChunkCount = 6

	minAutoChunkSize = 1024
	maxAutoChunkSize = 1073741824
)

// WriteRequest represents a blob to be written alone with metadata.
type WriteRequest struct {
	Id       string
	RollSize uint64
	Fields   Fields
	Blob     io.ReadCloser
}

// NewWrite creates a new WriteRequest to be used with Fixity.WriteRequest.
func NewWrite(id string, rc io.ReadCloser, f ...Field) *WriteRequest {
	return &WriteRequest{
		Id:       id,
		RollSize: DefaultAverageChunkSize,
		Fields:   f,
		Blob:     rc,
	}
}

func (req *WriteRequest) setRollSizeWithMin(roll uint64) {
	if roll < minAutoChunkSize {
		roll = minAutoChunkSize
	}
	if roll > maxAutoChunkSize {
		roll = maxAutoChunkSize
	}
	req.RollSize = roll
}

func (req *WriteRequest) SetRollFromBytes(b []byte) {
	req.setRollSizeWithMin(uint64(len(b)) / autoChunkCount)
}

func (req *WriteRequest) SetRollFromFileInfo(fi os.FileInfo) {
	req.setRollSizeWithMin(uint64(fi.Size()) / autoChunkCount)
}
