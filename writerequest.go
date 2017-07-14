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
	// This value is 1024^2*1 bytes, 1MiB, chosen as a compromise for the
	// small and large files.
	//
	// Ideally it should strike a balance between too many chunks written while
	// still mitigating data duplication caused by frequent changes to files
	// and written data.
	//
	// If written data is commonly editing the first X bytes, this value should
	// be overridden in the WriteRequest itself.
	DefaultAverageChunkSize = medChunkSize

	// smallChunkSize is the size SetChunkSizeFromX will use for small data.
	//
	// This value is 1024*100 bytes, 100 KiB, chosen as a balance for very
	// small files being able to be chunked but still being well over the
	// hard limits of the Restic "average bits". At this time, that value is
	// 14bits (12KiB). Reference for this value can be found here:
	//
	//    https://github.com/leeola/fixity/blob/master/chunkers/restic/restic.go#L12
	//
	// Note that that using this smallChunkSize is decided via the smallCutoff,
	// below.
	smallChunkSize uint64 = 102400

	// medChunkSize is the size SetChunkSizeFromX will use for medium data.
	//
	// This value is 1024^2 bytes, 1MiB, chosen as a balance for medium data.
	//
	// Note that using this medChunkSize is decided via the medCutoff, below.
	medChunkSize uint64 = 1048576

	// medChunkSize is the size SetChunkSizeFromX will use for large data.
	//
	// This value is 1024^2*4 bytes, 4MiB, chosen as a balance between total
	// hashes created and size of chunks being sent over the wire for large
	// data. Too big of a chunk would cause data being sent from one node
	// to another to lag. This value should choose a middleground between
	// total hashes and node transfer times.
	//
	// Note that largeChunkSize is chosen only if the total bytes are larger
	// the medCutoff.
	largeChunkSize uint64 = 4194304

	// smallCutoff decides if the smallChunkSize is to be used or not.
	//
	// If the total bytes given to AssertChunkSize() is below this value,
	// the smallChunkSize will be used.
	smallCutoff = smallChunkSize * 10

	// medCutoff decides if the medChunkSize is to be used or not.
	//
	// If the total bytes given to AssertChunkSize() is below this value,
	// the medChunkSize will be used.
	medCutoff = medChunkSize * 200
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

// SetChunkSizeFromBytes sets the chunkSize from the given byte slice.
//
// The chunksize value is decided via AssertChunkSize().
func (req *WriteRequest) SetChunkSizeFromBytes(b []byte) {
	req.AverageChunkSize = AssertChunkSize(uint64(len(b)))
}

// SetChunkSizeFromFileInfo sets the chunkSize from the given FileInfo.
//
// The chunksize value is decided via AssertChunkSize().
func (req *WriteRequest) SetChunkSizeFromFileInfo(fi os.FileInfo) {
	req.AverageChunkSize = AssertChunkSize(uint64(fi.Size()))
}

// AssertChunkSize chooses a chunksize based on the byte size of the input.
//
// It uses a tiered approach, decided by the constant values smallCutoff and
// medCutoff. In short:
//
// - If totalBytes is smaller than 1MiB, a chunkSize of 100KiB is used.
// - If totalBytes is smaller than 200MiB, a chunkSize of 1MiB is used.
// - If totalBytes is larger than 200MiB, a chunkSize of 4MiB is used.
func AssertChunkSize(totalBytes uint64) (chunkSize uint64) {
	switch {
	case totalBytes < smallCutoff:
		return smallChunkSize
	case totalBytes < medCutoff:
		return medChunkSize
	default:
		return largeChunkSize
	}
}
