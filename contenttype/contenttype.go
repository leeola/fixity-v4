package contenttype

import "io"

// ContentType processes incoming data for a specific type with supplied metadata.
//
// This allows the caller of the /upload/:type api to inform Kala of metadata about
// the data being uploaded. The Kala upload plugin (such as a jpeg or mp3 plugin)
// will do the bulk of the actual processing, but this interface ensures that
// arbitrary metadata that already exists is not lost on upload.
//
// Does it need to be parse exif data? Does it need to parse mp3 tags? Etc.
//
// The ContentType is responsible for writing raw blobs as needed. If multipart
// or blob chunks need to be written, it is responsible for doing so!
type ContentType interface {
	// StoreContent stores and indexes content with the given meta changes.
	StoreContent(io.ReadCloser, Version, Changes) ([]string, error)

	// StoreMeta stores and indexes the version and meta changes.
	StoreMeta(Version, Changes) ([]string, error)

	// UnmarshalMeta indexes already stored version and given metadata bytes.
	//
	// The purpose of this is to restore an index for the given Metadata.
	UnmarshalMeta([]byte) (interface{}, error)
}
