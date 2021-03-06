package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/reader/blobreader"
	"github.com/leeola/fixity/value"
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

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	if useStdin {
		return writeReadCloser(clictx, s, ioutil.NopCloser(os.Stdin), id)
	}

	for _, filename := range filenames {
		if err := writeFile(clictx, s, id, filename); err != nil {
			return fmt.Errorf("writereadcloser %q: %v", filename, err)
		}
	}

	return nil
}

func writeFile(clictx *cli.Context, s store, id, filename string) error {
	if id == "" {
		paths := []string{"files"}
		if dir := filepath.Base(filepath.Dir(filename)); dir != "" {
			paths = append(paths, dir)
		}
		paths = append(paths, filepath.Base(filename))
		id = filepath.Join(paths...)
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("openfile %q: %v", filename, err)
	}
	defer f.Close()

	if err := writeReadCloser(clictx, s, f, id); err != nil {
		return fmt.Errorf("writereadcloser %q: %v", filename, err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	return nil
}

func writeReadCloser(clictx *cli.Context, s store, r io.Reader, id string) error {
	preview := clictx.Bool("preview")
	allowUnsafe := clictx.Bool("allow-unsafe")

	if id == "" {
		return errors.New("id must be defined if it cannot be inferred")
	}

	var values fixity.Values
	for _, kv := range clictx.StringSlice("kv") {
		if values == nil {
			values = fixity.Values{}
		}
		k, v, err := splitKV(kv)
		if err != nil {
			return err // no wrap helper err
		}
		values[k] = value.String(v)
	}

	hashes, err := s.Write(context.Background(), id, values, r)
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

	return nil
}

func splitKV(kv string) (string, string, error) {
	split := strings.SplitN(kv, "=", 2)
	if len(split) != 2 {
		return "", "", errors.New("invalid kv format, requires key=value pair")
	}

	return split[0], split[1], nil
}
