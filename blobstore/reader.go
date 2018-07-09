package blobstore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/leeola/fixity"
)

type Reader interface {
	Read(context.Context, fixity.Ref) (io.ReadCloser, error)
}

func ReadAndUnmarshal(ctx context.Context, r Reader, ref fixity.Ref, v interface{}) error {
	rc, err := r.Read(ctx, ref)
	if err != nil {
		return fmt.Errorf("blobstore read: %v", err)
	}
	defer rc.Close()

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("readall: %v", err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}
