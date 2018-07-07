package main

import (
	"context"
	"errors"
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

	id := clictx.String("id")

	filenames := clictx.Args()
	filenamesLen := len(filenames)
	useFiles := filenamesLen > 0

	if filenamesLen > 1 && id != "" {
		return errors.New("cannot write multiple files to a single id")
	}
	if !useFiles && !useStdin {
		return errors.New("missing files or stdin to write")
	}
	if useFiles && id == "" {
		id = filenames[0]
	}

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	if useStdin {
		return writeReadCloser(clictx, s, ioutil.NopCloser(os.Stdin), id)
	}

	for _, filename := range filenames {
		f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("openfile %q: %v", filename, err)
		}

		// func closes
		if err := writeReadCloser(clictx, s, f, id); err != nil {
			return fmt.Errorf("writereadcloser %q: %v", filename, err)
		}
	}

	return nil
}

func writeReadCloser(clictx *cli.Context, s store, rc io.ReadCloser, id string) error {
	defer rc.Close()

	preview := clictx.Bool("preview")
	allowUnsafe := clictx.Bool("allow-unsafe")

	if id == "" {
		return errors.New("id must be defined if it cannot be inferred")
	}

	hashes, err := s.Write(context.Background(), id, nil, rc)
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
