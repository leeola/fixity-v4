package file

import (
	"io"
	"io/ioutil"
	"time"

	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/upload"
)

func FileUpload(s store.Store) upload.UploadFunc {
	return func(rc io.ReadCloser, m upload.Metadata) ([]string, error) {
		defer rc.Close()
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, errors.Stack(err)
		}

		now := time.Now()

		var contentHashes []string
		// TODO(leeola): in the future, split the file into deduped chunks here.
		h, err := store.WriteContent(s, store.Content{
			Type:    store.ContentType,
			Content: b,
		})
		if err != nil {
			return nil, err
		}
		contentHashes = append(contentHashes, h)

		h, err = store.WriteMultiPart(s, store.MultiPart{
			Type:      store.MultiPartType,
			CreatedAt: now,
			Parts:     contentHashes,
		})
		if err != nil {
			return nil, err
		}

		return append([]string{h}, contentHashes...), nil
	}
}
