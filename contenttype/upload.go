package contenttype

import (
	"io"

	"github.com/leeola/kala/store"
)

// Uploader processes incoming data for a specific type with supplied metadata.
//
// This allows the caller of the /upload/:type api to inform Kala of metadata about
// the data being uploaded. The Kala upload plugin (such as a jpeg or mp3 plugin)
// will do the bulk of the actual processing, but this interface ensures that
// arbitrary metadata that already exists is not lost on upload.
//
// Does it need to be parse exif data? Does it need to parse mp3 tags? Etc.
//
// The uploader is responsible for writing raw blobs as needed. If multipart or
// permanode chunks need to be written, it is responsible for doing so!
type Uploader interface {
	Upload(io.ReadCloser, store.MetaChanges) ([]string, error)
}

type Index interface {
}

// Restorer is a general purpose interface for restoring data to it's original state.
//
// This is not frequently implemented, as the use cases are not always applicable,
// but the easiest example is restoring posix metadata for a file back to the file.
//
// NOTE: This is not needed anywhere currently. Hence the lack of any methods ;)
type Restorer interface {
}

// UploadFunc implements Upload for a single function.
type UploadFunc (func(io.ReadCloser, store.MetaChanges) ([]string, error))

func (f UploadFunc) Upload(rc io.ReadCloser, c store.MetaChanges) ([]string, error) {
	return f(rc, c)
}
