package upload

import (
	"io"

	"github.com/leeola/kala/store"
)

// MetaChanges is a map of metadata changes to write along with the upload data.
//
// It is up to the uploader to decide how this metadata is stored within the
// Meta blob, the caller should make no assumptions. See specific Upload
// implementation documentation as needed.
//type MetaChanges map[string]string

// Upload processes incoming data for a specific type with supplied metadata.
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
//
// TODO(leeola): refactor Upload to a ContentType interface with an Upload and
// Download method, primarily shuttling content specific metadata into and
// out of the store.
type Upload interface {
	Upload(io.ReadCloser, store.MetaChanges) ([]string, error)
}

// UploadFunc implements Upload for a single function.
type UploadFunc (func(io.ReadCloser, store.MetaChanges) ([]string, error))

func (f UploadFunc) Upload(rc io.ReadCloser, c store.MetaChanges) ([]string, error) {
	return f(rc, c)
}
