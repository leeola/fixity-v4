package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/reader/blobreader"
	"github.com/urfave/cli"
)

func WriteCmd(clictx *cli.Context) error {
	useStdin := clictx.Bool("stdin")

	var r io.Reader
	switch {
	case useStdin:
		r = os.Stdin
	}

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	preview := clictx.Bool("preview")
	allowUnsafe := clictx.Bool("allow-unsafe")

	id := "foo"

	hashes, err := s.Write(context.Background(), id, nil, r)
	if err != nil {
		return fmt.Errorf("write: %v", err)
	}

	for _, h := range hashes {
		fmt.Println(h)

		if preview {
			if err := previewBlob(context.Background(), s, h, allowUnsafe); err != nil {
				return fmt.Errorf("previewblob: %v", err)
			}
		}
	}

	return nil
}

func previewBlob(ctx context.Context, s store, ref fixity.Ref, notSafe bool) error {
	rc, err := s.Blob(ctx, ref)
	if err != nil {
		return fmt.Errorf("blob: %v", err)
	}
	defer rc.Close()

	r, bt, err := blobreader.BlobType(rc)
	if err != nil {
		return fmt.Errorf("blobtype: %v", err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("readall: %v", err)
	}

	switch {
	case bt != fixity.BlobTypeSchemaless:
		if err := printJsonBytes(os.Stdout, b); err != nil {
			return fmt.Errorf("printjsonbytes: %v", err)
		}
	case notSafe:
		fmt.Println(string(b))
	}

	// newline
	fmt.Println("")

	return nil
}
