package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/leeola/fixity"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

func ReadCmd(clictx *cli.Context) error {
	if len(clictx.Args()) != 1 {
		return errors.New("missing mutation reference argument")
	}

	color.NoColor = clictx.Bool("no-stderr-color")

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	ref := fixity.Ref(clictx.Args().First())
	mutation, values, r, err := s.ReadRef(context.Background(), ref)
	if err != nil {
		return fmt.Errorf("read %q: %v", ref, err)
	}

	if !clictx.Bool("no-mutation") {
		fmt.Fprintln(os.Stderr, "mutation:")
		if err := printAsJSON(os.Stderr, mutation); err != nil {
			return fmt.Errorf("print mutation: %v", err)
		}
	}

	if !clictx.Bool("no-values") && values != nil {
		fmt.Fprintln(os.Stderr, "values:")
		if err := printAsJSON(os.Stderr, values); err != nil {
			return fmt.Errorf("print mutation: %v", err)
		}
	}

	filename := clictx.String("filename")
	if filename != "" {
		fmt.Fprintln(os.Stderr, "data written to:", filename)
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
	} else {

		var redirectedText string
		redirected := os.Getenv("TERM") == "dumb" ||
			(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
		if redirected {
			redirectedText = "redirected"
		}

		fmt.Fprintln(os.Stderr, "data:", redirectedText)

		if _, err := io.Copy(os.Stdout, r); err != nil {
			return fmt.Errorf("copy stdout: %v", err)
		}

		if !redirected {
			fmt.Fprintln(os.Stderr)
		}
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
