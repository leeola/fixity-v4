package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func ReadCmd(clictx *cli.Context) error {
	if len(clictx.Args()) > 3 {
		return errors.New("too many args")
	}

	filename := clictx.Args().Get(1)
	if filename == "" {
		return fmt.Errorf("missing filename arg")
	}

	color.NoColor = clictx.Bool("no-stderr-color")

	dataMsg := "data: written to " + filename
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("openfile %q: %v", filename, err)
	}
	defer f.Close()

	if err := readToWriters(clictx, f, os.Stderr, dataMsg); err != nil {
		// no wrap above helper errs
		return err
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %v", err)
	}

	return nil
}

// printAsJSON marshalls the given struct to json to print it with the
// same highlighter as blob printing.
func printAsJSON(out io.Writer, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal: %v", err)
	}

	return printJsonBytes(out, b)
}
