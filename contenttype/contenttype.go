package contenttype

import (
	"io"

	"github.com/leeola/kala/store"
)

type Importer interface {
	Import(map[string]string) (io.ReadCloser, store.MetaChanges, error)
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
	StoreContent(io.ReadCloser, store.MetaChanges) ([]string, error)

	// Meta applies just metadata changes.
	//
	// Note that this just changes metadata, but it can change the content that
	// the metadata points to.
	Meta(store.MetaChanges) ([]string, error)
}

// Exporter is a general purpose interface for restoring data to it's original state.
//
// This is not frequently implemented, as the use cases are not always applicable,
// but the easiest example is restoring posix metadata for a file back to the file.
type Exporter interface {
	Export(io.ReadCloser, string) error
}
