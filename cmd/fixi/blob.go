package main

import (
	"context"
	"fmt"

	"github.com/leeola/fixity"
	"github.com/urfave/cli"
)

func BlobCmd(clictx *cli.Context) error {
	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	for _, sRef := range clictx.Args() {
		ref := fixity.Ref(sRef)
		if err := printBlob(context.Background(), s, ref); err != nil {
			return fmt.Errorf("printblob %q: %v", ref, err)
		}
	}

	return nil
}

type store interface {
	Blob(ctx context.Context, ref fixity.Ref) (fixity.BlobReadCloser, error)
}

func printBlob(ctx context.Context, s store, ref fixity.Ref) error {
	rc, err := s.Blob(context.Background(), ref)
	if err != nil {
		return fmt.Errorf("blob: %v", err)
	}

	bt, err := rc.BlobType()
	if err != nil {
		return fmt.Errorf("blobtype: %v", err)
	}

	fmt.Println(ref, ",", bt)

	return nil
}
