package file

import (
	"io"

	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
	"github.com/leeola/kala/store/roller/camli"
	"github.com/leeola/kala/upload"
)

func FileUpload(s store.Store) upload.UploadFunc {
	return func(rc io.ReadCloser, m upload.Metadata) ([]string, error) {
		if rc == nil {
			return nil, errors.New("missing ReadCloser")
		}
		defer rc.Close()

		roller, err := camli.New(rc)
		if err != nil {
			return nil, errors.Stack(err)
		}

		hashes, err := store.WriteContentRoller(s, roller)
		if err != nil {
			return nil, errors.Stack(err)
		}

		h, err := store.WriteMultiPart(s, store.MultiPart{
			Parts: hashes,
		})
		if err != nil {
			return nil, errors.Stack(err)
		}

		return append([]string{h}, hashes...), nil
	}
}
