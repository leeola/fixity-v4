package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/leeola/fixity"
	"github.com/urfave/cli"
)

func ReadCmd(clictx *cli.Context) error {
	if len(clictx.Args()) != 1 {
		return errors.New("missing mutation reference argument")
	}

	filename := clictx.String("filename")
	if filename == "" {
		return fmt.Errorf("filename currently required")
	}

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	ref := fixity.Ref(clictx.Args().First())
	r, err := s.Read(context.Background(), ref)
	if err != nil {
		return fmt.Errorf("read %q: %v", ref, err)
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("openfile %q: %v", filename, err)
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("copy file: %v", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	return nil
}
