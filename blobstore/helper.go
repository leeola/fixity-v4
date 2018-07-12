package blobstore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/leeola/fixity"
)

func ReadAndUnmarshal(ctx context.Context, r fixity.BlobReader, ref fixity.Ref, v interface{}) error {
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
