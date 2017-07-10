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

	// autoChunkCount is the number of chunks SetChunkFromBytes/etc will set to.
	//
	// Eg, if SetChunkFromBytes is given a byte array that is 1,048,576 bytes,
	// then it will divide 1,048,576 by eg, 10 to set the rollSize to
	// 104,857, effectively rolling the bytes into 10 parts when it's written.
	autoChunkCount = 6

	// 12KiB so that small files/data won't churn with a ton of chunks.
	minAutoChunkSize = 12288
	// 4MiB because chunks too large become hard to store/send/etc. Each
	// chunk should be reasonbly fast to send between nodes.
	maxAutoChunkSize = DefaultAverageChunkSize
)

// WriteRequest represents a blob to be written alone with metadata.
type WriteRequest struct {
	// Id is the Content Id used to associate this write with previous writes.
	//
	// This association effectively mutates the write from the previous content
	// with the same Id. Important values are inhereted, such as ChunkSize,
	// ensuring consistent writes.
	Id string `json:"id,omitempty"`

	// AverageChunkSize is the average bytes that each chunk is aimed to be.
	//
	// Chunks are separated by Cotent Defined Chunks (CDC) and this value
	// allows mutations of this blob to use the same ChunkSize with each
	// version. This ensures the chunks are chunk'd by the CDC algorithm
	// with the same spacing.
	//
	// Note that the algorithm is decided by the fixity.Store.
	AverageChunkSize uint64 `json:"averageChunkSize,omitempty"`

	// IgnoreDuplicateBlob will ignore this write if the blob exists for the id.
	//
	// An existing blob is determined to be one which *exists* on the chain,
	// not specifically if it is the latest Content for it's id.
	//
	// This flag is used to ensure data exists on the blockchain, without
	// causing any form of versioning contention. Duplicate blobs are allowed
	// for different Ids.
	IgnoreDuplicateBlob bool

	// Fields are the indexable fields that the data will be indexed with.
	//
	// This is often metadata like filename, unix permissions, etc. It can
	// also be used to retrieve this data later, as the Fields is stored
	// on the fixity.Content object within the datastore.
	Fields Fields `json:"fields,omitempty"`

	// Blob is the actual data that is being written, and is required.
	//
	// This will be used to build a fixity.Blob type.
	Blob io.ReadCloser `json:"-"`
}

// NewWrite creates a new WriteRequest to be used with Fixity.WriteRequest.
func NewWrite(id string, rc io.ReadCloser, f ...Field) *WriteRequest {
	return &WriteRequest{
		Id:               id,
		AverageChunkSize: DefaultAverageChunkSize,
		Fields:           f,
		Blob:             rc,
	}
}

func (req *WriteRequest) setChunkSize(averageChunkSize uint64) {
	if averageChunkSize < minAutoChunkSize {
		averageChunkSize = minAutoChunkSize
	}
	if averageChunkSize > maxAutoChunkSize {
		averageChunkSize = maxAutoChunkSize
	}
	req.AverageChunkSize = averageChunkSize
}

func (req *WriteRequest) SetChunkSizeFromBytes(b []byte) {
	req.setChunkSize(uint64(len(b)) / autoChunkCount)
}

func (req *WriteRequest) SetChunkSizeFromFileInfo(fi os.FileInfo) {
	req.setChunkSize(uint64(fi.Size()) / autoChunkCount)
}
