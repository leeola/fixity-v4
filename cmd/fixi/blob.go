package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/leeola/fixity"
	"github.com/nwidger/jsoncolor"
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
	rc, err := s.Blob(ctx, ref)
	if err != nil {
		return fmt.Errorf("blob: %v", err)
	}
	defer rc.Close()

	// Disabled currently.
	// bt, err := rc.BlobType()
	// if err != nil {
	// 	return fmt.Errorf("blobtype: %v", err)
	// }

	b, err := ioutil.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("readall: %v", err)
	}

	if err := printJsonBytes(os.Stdout, b); err != nil {
		return fmt.Errorf("printjsonbytes: %v", err)
	}

	return nil
}

func printJsonBytes(out io.Writer, b []byte) error {
	f := jsoncolor.NewFormatter()

	f.SpaceColor = color.New(color.FgRed, color.Bold)
	f.CommaColor = color.New(color.FgWhite, color.Bold)
	f.ColonColor = color.New(color.FgBlue)
	f.ObjectColor = color.New(color.FgBlue, color.Bold)
	f.ArrayColor = color.New(color.FgWhite)
	f.FieldColor = color.New(color.FgGreen)
	f.StringColor = color.New(color.FgBlack, color.Bold)
	f.TrueColor = color.New(color.FgWhite, color.Bold)
	f.FalseColor = color.New(color.FgRed)
	f.NumberColor = color.New(color.FgWhite)
	f.NullColor = color.New(color.FgWhite, color.Bold)

	prettyJson := bytes.Buffer{}
	if err := f.Format(&prettyJson, b); err != nil {
		return err
	}

	if _, err := io.Copy(out, &prettyJson); err != nil {
		return err
	}

	fmt.Print("\n")
	return nil
}
