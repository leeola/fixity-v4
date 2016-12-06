package upload

import "io"

// Metadata is a map of metadata to write along with the upload data.
//
// It is up to the uploader to decide how this metadata is stored with the
// content, the caller should make no assumptions. See specific Upload
// implementation documentation as needed.
type Metadata map[string]interface{}

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
type Upload interface {
	Upload(io.ReadCloser, Metadata) ([]string, error)
}

// UploadFunc implements Upload for a single function.
type UploadFunc (func(io.ReadCloser, Metadata) ([]string, error))

func (f UploadFunc) Upload(rc io.ReadCloser, m Metadata) ([]string, error) {
	return f(rc, m)
}
