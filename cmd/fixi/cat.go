package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/leeola/fixity"
	"github.com/mattn/go-isatty"
	"github.com/urfave/cli"
)

func CatCmd(clictx *cli.Context) error {
	if len(clictx.Args()) > 1 {
		return errors.New("too many args")
	}

	wout, werr := os.Stdout, os.Stderr

	var redirectedText string
	redirected := os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
	if redirected {
		redirectedText = "redirected"
	}

	dataMsg := "data: " + redirectedText

	if err := readToWriters(clictx, wout, werr, dataMsg); err != nil {
		// no wrap above helper errs
		return err
	}

	if !redirected {
		fmt.Fprintln(os.Stderr)
	}

	return nil
}

func readToWriters(clictx *cli.Context, wout, werr io.Writer, dataMsg string) error {
	idOrRef := clictx.Args().Get(0)
	if idOrRef == "" {
		return fmt.Errorf("missing id or ref arg")
	}

	color.NoColor = clictx.Bool("no-stderr-color")

	s, err := storeFromCli(clictx)
	if err != nil {
		// no wrap above helper errs
		return err
	}

	var (
		mutation fixity.Mutation
		values   fixity.Values
		r        fixity.Reader
	)

	if clictx.Bool("ref") {
		ref := fixity.Ref(idOrRef)
		mutation, values, r, err = s.ReadRef(context.Background(), ref)
		if err != nil {
			return fmt.Errorf("read %q: %v", ref, err)
		}
	} else {
		id := idOrRef
		mutation, values, r, err = s.Read(context.Background(), id)
		if err != nil {
			return fmt.Errorf("read %q: %v", id, err)
		}
	}

	if !clictx.Bool("no-mutation") {
		fmt.Fprintln(werr, "mutation:")
		if err := printAsJSON(werr, mutation); err != nil {
			return fmt.Errorf("print mutation: %v", err)
		}
	}

	if !clictx.Bool("no-values") && values != nil {
		fmt.Fprintln(werr, "values:")
		if err := printAsJSON(werr, values); err != nil {
			return fmt.Errorf("print mutation: %v", err)
		}
	}

	fmt.Fprintln(werr, dataMsg)
	if _, err := io.Copy(wout, r); err != nil {
		return fmt.Errorf("copy wout: %v", err)
	}

	return nil
}
