package contenttype

import "io"

type Importer interface {
	Import(map[string]string) (io.ReadCloser, MetaChanges, error)
}

// ContentStorer processes incoming data for a specific type with supplied metadata.
//
// This allows the caller of the /upload/:type api to inform Kala of metadata about
// the data being uploaded. The Kala upload plugin (such as a jpeg or mp3 plugin)
// will do the bulk of the actual processing, but this interface ensures that
// arbitrary metadata that already exists is not lost on upload.
//
// Does it need to be parse exif data? Does it need to parse mp3 tags? Etc.
//
// The ContentStorer is responsible for writing raw blobs as needed. If multipart
// or blob chunks need to be written, it is responsible for doing so!
type ContentStorer interface {
	// StoreContent stores content with the given meta changes.
	//
	// Optionally, the metadata can be passed in as a byte array. This is to allow
	// the caller of this interface to read the metadata from the store while
	// avoiding a double read on the metadata.
	//
	// The implementor *must* handle both an empty byte array, and a populated array.
	StoreContent(io.ReadCloser, []byte, MetaChanges) ([]string, error)

	// Meta applies just metadata changes.
	//
	// Optionally, the metadata can be passed in as a byte array. This is to allow
	// the caller of this interface to read the metadata from the store while
	// avoiding a double read on the metadata.
	//
	// The implementor *must* handle both an empty byte array, and a populated array.
	Meta([]byte, MetaChanges) ([]string, error)
}

// Exporter is a general purpose interface for restoring data to it's original state.
//
// This is not frequently implemented, as the use cases are not always applicable,
// but the easiest example is restoring posix metadata for a file back to the file.
type Exporter interface {
	Export(io.ReadCloser, string) error
}
